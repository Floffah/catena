import { auth } from "@clerk/nextjs/server";
import { cache } from "react";

import { setAuthToken } from "@/lib/api";

export const authenticateApiClient = cache(async () => {
    const { getToken, isAuthenticated } = await auth();

    if (isAuthenticated) {
        setAuthToken(await getToken());
    }
});
