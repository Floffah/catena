import { cache } from "react";

import { serverGetApiClient } from "@/lib/server/auth";

export const serverGetGitAccessTokens = cache(async () => {
    const apiClient = await serverGetApiClient();

    const res = await apiClient.GET("/v1/git-access-tokens");

    return res.data;
});
