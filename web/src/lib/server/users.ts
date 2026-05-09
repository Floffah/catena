import { cache } from "react";

import { apiFetch } from "@/lib/api";
import { authenticateApiClient } from "@/lib/server/auth";

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
