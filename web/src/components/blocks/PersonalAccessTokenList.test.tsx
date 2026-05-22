import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, test } from "bun:test";
import { HttpResponse, http } from "msw";
import { act } from "react";

import PersonalAccessTokenList from "@/components/blocks/PersonalAccessTokenList";
import { renderWithQueryClient } from "@/test/render";
import { server } from "@/test/server";
import { SchemaGitAccessToken } from "@/types/api";

const token: SchemaGitAccessToken = {
    createdAt: "2026-05-22T00:00:00Z",
    expiresAt: null,
    id: "019deb10-dafc-743f-8cfc-289a80c13af1",
    lastUsedAt: null,
    name: "Local development laptop",
    revokedAt: null,
    scopes: ["repo:read", "repo:write"],
    tokenPrefix: "ctn_pat_abc12345",
    updatedAt: "2026-05-22T00:00:00Z",
};

describe("PersonalAccessTokenList", () => {
    test("shows an empty state when there are no active tokens", async () => {
        server.use(
            http.get("http://catena.test/v1/git-access-tokens", () =>
                HttpResponse.json([]),
            ),
        );

        await act(async () => {
            renderWithQueryClient(<PersonalAccessTokenList />);
        });

        expect(
            await screen.findByText("No active personal access tokens"),
        ).toBeDefined();
        expect(
            screen.getByText(
                "Create a token when you need to authenticate Git operations over HTTPS.",
            ),
        ).toBeDefined();
    });

    test("revokes a token and refetches the list", async () => {
        const deletedTokenIds: string[] = [];
        let tokens = [token];

        server.use(
            http.get("http://catena.test/v1/git-access-tokens", () =>
                HttpResponse.json(tokens),
            ),
            http.delete(
                "http://catena.test/v1/git-access-tokens/:id",
                ({ params }) => {
                    deletedTokenIds.push(String(params.id));
                    tokens = [];

                    return new Response(null, {
                        status: 204,
                    });
                },
            ),
        );

        await act(async () => {
            renderWithQueryClient(<PersonalAccessTokenList />);
        });

        expect(
            await screen.findByText("Local development laptop"),
        ).toBeDefined();
        expect(
            screen.getByText("Starts with ctn_pat_abc12345..."),
        ).toBeDefined();
        expect(screen.getByText("Never expires")).toBeDefined();

        await userEvent.click(
            screen.getByRole("button", {
                name: "Revoke Local development laptop",
            }),
        );

        await waitFor(() => {
            expect(deletedTokenIds).toEqual([
                "019deb10-dafc-743f-8cfc-289a80c13af1",
            ]);
        });

        expect(
            await screen.findByText("No active personal access tokens"),
        ).toBeDefined();
    });
});
