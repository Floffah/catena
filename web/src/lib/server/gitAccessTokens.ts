import { cache } from "react";

import { apiFetch } from "@/lib/api";
import { authenticateApiClient } from "@/lib/server/auth";

export const serverGetGitAccessTokens = cache(async () => {
    await authenticateApiClient();

    const res = await apiFetch.GET("/v1/git-access-tokens");

    return res.data;
});
