import createFetchClient, { Middleware } from "openapi-fetch";
import createQueryClient from "openapi-react-query";

import { useCatenaAuth } from "@/providers/AuthProvider";
import { paths } from "@/types/api";

export const apiFetch = createFetchClient<paths>({
    baseUrl:
        process.env.CATENA_DIRECT_INSTANCE_URL ||
        process.env.NEXT_PUBLIC_CATENA_INSTANCE_URL ||
        "http://localhost:8080",
});

export const $api = createQueryClient(apiFetch);

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
