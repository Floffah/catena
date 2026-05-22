import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, test } from "bun:test";
import { HttpResponse, http } from "msw";

import { renderWithQueryClient } from "@/test/render";
import { server } from "@/test/server";

describe("Create personal access token page", () => {
    test("creates a token and lets the user copy the raw secret once", async () => {
        const clipboardWrites: string[] = [];
        const requestBodies: unknown[] = [];

        Object.defineProperty(navigator, "clipboard", {
            configurable: true,
            value: {
                writeText: async (value: string) => {
                    clipboardWrites.push(value);
                },
            },
        });

        server.use(
            http.post(
                "http://catena.test/v1/git-access-tokens",
                async ({ request }) => {
                    requestBodies.push(await request.json());

                    return HttpResponse.json({
                        accessToken: {
                            createdAt: "2026-05-22T00:00:00Z",
                            expiresAt: null,
                            id: "019deb10-dafc-743f-8cfc-289a80c13af1",
                            lastUsedAt: null,
                            name: "Local development laptop",
                            revokedAt: null,
                            scopes: ["repo:read", "repo:write"],
                            tokenPrefix: "ctn_pat_secret",
                            updatedAt: "2026-05-22T00:00:00Z",
                        },
                        token: "ctn_pat_secret-value",
                    });
                },
            ),
        );

        const Page = await import("./page").then((mod) => mod.default);

        renderWithQueryClient(<Page />);

        await userEvent.type(
            screen.getByLabelText("Token Name"),
            "Local development laptop",
        );
        await userEvent.click(
            screen.getByRole("button", {
                name: "Create Token",
            }),
        );

        await waitFor(() => {
            expect(requestBodies).toEqual([
                {
                    name: "Local development laptop",
                },
            ]);
        });

        expect(await screen.findByText("Copy your token now")).toBeDefined();
        expect(
            screen.getByText(
                "This is the only time we will show this token. Copy it now before continuing.",
            ),
        ).toBeDefined();
        expect(screen.getByDisplayValue("ctn_pat_secret-value")).toBeDefined();

        await userEvent.click(
            screen.getByRole("button", {
                name: "Copy token",
            }),
        );

        expect(clipboardWrites).toEqual(["ctn_pat_secret-value"]);
    });
});
