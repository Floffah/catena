import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, test } from "bun:test";
import { HttpResponse, http } from "msw";

import { routerPushCalls, routerRefreshCalls } from "@/test/navigation";
import { renderWithQueryClient } from "@/test/render";
import { server } from "@/test/server";

describe("CreateIssueForm", () => {
    test("creates an issue and navigates to it", async () => {
        const requestBodies: unknown[] = [];

        server.use(
            http.post(
                "http://catena.test/v1/repositories/:owner/:repository/issues",
                async ({ request }) => {
                    requestBodies.push(await request.json());

                    return HttpResponse.json(
                        {
                            authorId: "019deb10-dafc-743f-8cfc-289a80c13af1",
                            body: "This is the issue body.",
                            createdAt: "2026-05-22T00:00:00Z",
                            id: "019deb10-dafc-743f-8cfc-289a80c13af2",
                            kind: "issue",
                            lastActivityAt: "2026-05-22T00:00:00Z",
                            number: 3,
                            reference: "I-3",
                            repositoryId:
                                "019deb10-dafc-743f-8cfc-289a80c13af3",
                            status: "open",
                            title: "First issue",
                            updatedAt: "2026-05-22T00:00:00Z",
                        },
                        {
                            status: 201,
                        },
                    );
                },
            ),
        );

        const CreateIssueForm = await import("./CreateIssueForm").then(
            (mod) => mod.default,
        );

        renderWithQueryClient(
            <CreateIssueForm ownerName="floffah" repoName="catena" />,
        );

        await userEvent.type(screen.getByLabelText("Title"), "  First issue  ");
        await userEvent.type(
            screen.getByLabelText("Body"),
            "This is the issue body.",
        );
        await userEvent.click(
            screen.getByRole("button", {
                name: "Create issue",
            }),
        );

        await waitFor(() => {
            expect(requestBodies).toEqual([
                {
                    body: "This is the issue body.",
                    title: "First issue",
                },
            ]);
        });
        expect(routerPushCalls).toEqual(["/floffah/catena/issues/3"]);
        expect(routerRefreshCalls).toBe(1);
    });

    test("requires a title before creating an issue", async () => {
        const requestBodies: unknown[] = [];

        server.use(
            http.post(
                "http://catena.test/v1/repositories/:owner/:repository/issues",
                async ({ request }) => {
                    requestBodies.push(await request.json());

                    return HttpResponse.json({});
                },
            ),
        );

        const CreateIssueForm = await import("./CreateIssueForm").then(
            (mod) => mod.default,
        );

        renderWithQueryClient(
            <CreateIssueForm ownerName="floffah" repoName="catena" />,
        );

        await userEvent.click(
            screen.getByRole("button", {
                name: "Create issue",
            }),
        );

        expect(await screen.findByText("Title is required")).toBeDefined();
        expect(requestBodies).toEqual([]);
        expect(routerPushCalls).toEqual([]);
    });
});
