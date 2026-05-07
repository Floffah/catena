import { auth } from "@clerk/nextjs/server";
import { cache } from "react";

import { apiFetch, setAuthToken } from "@/lib/api";
import { SchemaRepositoryRefType } from "@/types/api";

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
    async (owner: string, repository: string, ref = "main", path = "/") => {
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
                        path,
                        ref,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetRepositoryTree = cache(
    async (owner: string, repository: string, ref = "main", path = "/") => {
        const { getToken, isAuthenticated } = await auth();

        if (isAuthenticated) {
            setAuthToken(await getToken());
        }

        const res = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}/tree",
            {
                params: {
                    path: {
                        owner,
                        repository,
                    },
                    query: {
                        ref,
                        path,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetRepositoryLatestCommit = cache(
    async (owner: string, repository: string, ref = "main", path = "/") => {
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
                    query: {
                        ref,
                        path,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverResolveRepositoryGitPath = cache(
    async (owner: string, repository: string, path: string) => {
        const { getToken, isAuthenticated } = await auth();

        if (isAuthenticated) {
            setAuthToken(await getToken());
        }

        const res = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}/git-path/resolve",
            {
                params: {
                    path: {
                        owner,
                        repository,
                    },
                    query: {
                        path,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverListRepositoryRefs = cache(
    async (
        owner: string,
        repository: string,
        type: SchemaRepositoryRefType = "branch",
    ) => {
        const { getToken, isAuthenticated } = await auth();

        if (isAuthenticated) {
            setAuthToken(await getToken());
        }

        const res = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}/refs",
            {
                params: {
                    path: {
                        owner,
                        repository,
                    },
                    query: {
                        type,
                    },
                },
            },
        );

        return res.data;
    },
);
