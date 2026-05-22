import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, test } from "bun:test";
import { HttpResponse, http } from "msw";

import { routerRefreshCalls } from "@/test/navigation";
import { renderWithQueryClient } from "@/test/render";
import { server } from "@/test/server";
import { SchemaRepository } from "@/types/api";

const repository: SchemaRepository = {
    createdAt: "2026-05-22T00:00:00Z",
    defaultBranch: "main",
    description: "Original description",
    id: "019deb10-dafc-743f-8cfc-289a80c13af1",
    name: "catena",
    ownerId: "019deb10-dafc-743f-8cfc-289a80c13af2",
    ownerName: "floffah",
    updatedAt: "2026-05-22T00:00:00Z",
    visibility: "private",
};

describe("RepositorySettingsForm", () => {
    test("updates repository settings", async () => {
        const requestBodies: unknown[] = [];

        server.use(
            http.patch(
                "http://catena.test/v1/repositories/:owner/:repository",
                async ({ request }) => {
                    requestBodies.push(await request.json());

                    return HttpResponse.json({
                        ...repository,
                        defaultBranch: "develop",
                        description: "Updated description",
                        visibility: "public",
                    });
                },
            ),
        );

        const RepositorySettingsForm =
            await import("./RepositorySettingsForm").then((mod) => mod.default);

        renderWithQueryClient(
            <RepositorySettingsForm
                branchNames={["main", "develop"]}
                repository={repository}
            />,
        );

        await userEvent.clear(screen.getByLabelText("Description (optional)"));
        await userEvent.type(
            screen.getByLabelText("Description (optional)"),
            "Updated description",
        );
        await userEvent.selectOptions(
            screen.getByLabelText("Default Branch"),
            "develop",
        );
        await userEvent.selectOptions(
            screen.getByLabelText("Visibility"),
            "public",
        );
        await userEvent.click(
            screen.getByRole("button", {
                name: "Save Changes",
            }),
        );

        await waitFor(() => {
            expect(requestBodies).toEqual([
                {
                    defaultBranch: "develop",
                    description: "Updated description",
                    visibility: "public",
                },
            ]);
        });
        expect(routerRefreshCalls).toBe(1);
        expect(
            await screen.findByText("Repository updated successfully!"),
        ).toBeDefined();
    });
});
