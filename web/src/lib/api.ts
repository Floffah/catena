import createFetchClient, { Middleware } from "openapi-fetch";
import createQueryClient from "openapi-react-query";

import { useCatenaAuth } from "@/providers/AuthProvider";

import { paths } from "../../types/api";

let token: string | null = null;

const authMiddleware: Middleware = {
    async onRequest({ request }) {
        if (token) {
            request.headers.set("Authorization", `Bearer ${token}`);
        }

        return request;
    },
};

export const apiFetch = createFetchClient<paths>({
    baseUrl:
        process.env.NEXT_PUBLIC_CATENA_INSTANCE_URL || "http://localhost:8080",
});
apiFetch.use(authMiddleware);

export const $api = createQueryClient(apiFetch);

export function setAuthToken(newToken: string | null) {
    token = newToken;
}

export const useAuthedQuery = ((
    method: unknown,
    url: unknown,
    init: unknown,
    options: unknown,
    queryClient: unknown,
) => {
    const auth = useCatenaAuth();

    return $api.useQuery(
        // @ts-expect-error - we can allow any here, type is guaranteed by the caller
        method,
        url,
        init,
        {
            ...(options ?? {}),
            enabled:
                auth.isAuthenticated &&
                (options as { enabled?: boolean })?.enabled !== false,
        },
        queryClient,
    );
}) as unknown as typeof $api.useQuery;
