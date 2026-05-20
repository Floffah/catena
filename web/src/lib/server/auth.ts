import { auth } from "@clerk/nextjs/server";
import createFetchClient from "openapi-fetch";
import { cache } from "react";

import { paths } from "@/types/api";

const baseUrl =
    process.env.CATENA_DIRECT_INSTANCE_URL ||
    process.env.NEXT_PUBLIC_CATENA_INSTANCE_URL ||
    "http://localhost:8080";

export const serverGetApiClient = cache(async () => {
    const { getToken, isAuthenticated } = await auth();
    const token = isAuthenticated ? await getToken() : null;

    const apiClient = createFetchClient<paths>({
        baseUrl,
    });

    if (token) {
        apiClient.use({
            async onRequest({ request }) {
                request.headers.set("Authorization", `Bearer ${token}`);

                return request;
            },
        });
    }

    return apiClient;
});
