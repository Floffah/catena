import { auth } from "@clerk/nextjs/server";
import { cache } from "react";

import { apiFetch, setAuthToken } from "@/lib/api";

export const serverGetRepository = cache(
    async (owner: string, repository: string) => {
        const { getToken, isAuthenticated } = await auth();

        if (isAuthenticated) {
            setAuthToken(await getToken());
        }

        const res = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}",
            {
                params: {
                    path: {
                        owner,
                        repository,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetRepositoryReadme = cache(
    async (owner: string, repository: string) => {
        const { getToken, isAuthenticated } = await auth();

        if (isAuthenticated) {
            setAuthToken(await getToken());
        }

        const res = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}/readme",
            {
                params: {
                    path: {
                        owner,
                        repository,
                    },
                    query: {
                        path: "/",
                        ref: "main",
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetRepositoryLatestCommit = cache(
    async (owner: string, repository: string) => {
        const { getToken, isAuthenticated } = await auth();

        if (isAuthenticated) {
            setAuthToken(await getToken());
        }

        const res = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}/latest-commit",
            {
                params: {
                    path: {
                        owner,
                        repository,
                    },
                },
            },
        );

        return res.data;
    },
);
