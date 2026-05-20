import { cache } from "react";

import { serverGetApiClient } from "@/lib/server/auth";
import { SchemaRepositoryRefType } from "@/types/api";

export const serverGetRepository = cache(
    async (ownerName: string, repoName: string) => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetRepositoryReadme = cache(
    async (ownerName: string, repoName: string, ref = "main", path = "/") => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/readme",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
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
    async (ownerName: string, repoName: string, ref = "main", path = "/") => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/tree",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
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
    async (ownerName: string, repoName: string, ref = "main", path = "/") => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/latest-commit",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
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
    async (ownerName: string, repoName: string, path: string) => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/git-path/resolve",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
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
        ownerName: string,
        repoName: string,
        type: SchemaRepositoryRefType = "branch",
    ) => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/refs",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
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

export const serverListRepositoryIssues = cache(
    async (ownerName: string, repoName: string) => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/issues",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetRepositoryIssue = cache(
    async (ownerName: string, repoName: string, number: number) => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET(
            "/v1/repositories/{owner}/{repository}/issues/{number}",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
                        number,
                    },
                },
            },
        );

        return res.data;
    },
);

export const serverGetCurrentRepositoryRef = cache(
    async (ownerName: string, repoName: string, path: string) => {
        const repo = await serverGetRepository(ownerName, repoName);

        if (repo) {
            const resolvedPath = await serverResolveRepositoryGitPath(
                ownerName,
                repoName,
                path,
            );

            return resolvedPath?.ref;
        }

        return undefined;
    },
);
