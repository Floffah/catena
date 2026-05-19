import { auth } from "@clerk/nextjs/server";
import { Middleware } from "openapi-fetch";
import { cache } from "react";

import { apiFetch } from "@/lib/api";

let serverToken: string | null = null;

const serverAuthMiddleware: Middleware = {
    async onRequest({ request }) {
        if (serverToken) {
            request.headers.set("Authorization", `Bearer ${serverToken}`);
        }

        return request;
    },
};
apiFetch.use(serverAuthMiddleware);

export const authenticateApiClient = cache(async () => {
    const { getToken, isAuthenticated } = await auth();

    if (isAuthenticated) {
        serverToken = await getToken();
    }
});
