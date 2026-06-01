import { cache } from "react";

import { serverGetApiClient } from "@/lib/server/auth";

export const serverGetCurrentUser = cache(async () => {
    const apiClient = await serverGetApiClient();

    const res = await apiClient.GET("/v1/user");

    return res.data;
});

export const serverGetUserForClerkID = cache(async (clerkID: string) => {
    const apiClient = await serverGetApiClient();

    const res = await apiClient.GET("/v1/users/clerk/{clerkUserId}", {
        params: {
            path: {
                clerkUserId: clerkID,
            },
        },
    });

    return res.data;
});

export const serverGetUserForName = cache(async (name: string) => {
    const apiClient = await serverGetApiClient();

    const res = await apiClient.GET("/v1/users/name/{name}", {
        params: {
            path: {
                name,
            },
        },
    });

    return res.data;
});

export const serverListFeaturedRepositoriesForUser = cache(
    async (name: string) => {
        const apiClient = await serverGetApiClient();

        const res = await apiClient.GET("/v1/users/name/{name}/repositories", {
            params: {
                path: {
                    name,
                },
                query: {
                    limit: 6,
                    sort: "featured",
                    visibility: "public",
                },
            },
        });

        return res.data;
    },
);
