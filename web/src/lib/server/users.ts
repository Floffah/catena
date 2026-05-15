import { cache } from "react";

import { apiFetch } from "@/lib/api";
import { authenticateApiClient } from "@/lib/server/auth";

export const serverGetCurrentUser = cache(async () => {
    await authenticateApiClient();

    const res = await apiFetch.GET("/v1/user");

    return res.data;
});

export const serverGetUserForClerkID = cache(async (clerkID: string) => {
    await authenticateApiClient();

    const res = await apiFetch.GET("/v1/users/clerk/{clerkUserId}", {
        params: {
            path: {
                clerkUserId: clerkID,
            },
        },
    });

    return res.data;
});

export const serverGetUserForName = cache(async (name: string) => {
    await authenticateApiClient();

    const res = await apiFetch.GET("/v1/users/name/{name}", {
        params: {
            path: {
                name,
            },
        },
    });

    return res.data;
});
