"use client";

import { $api, useAuthedQuery } from "@/lib/api";

export default function Page() {
    const health = useAuthedQuery("get", "/v1/healthz");

    return null;
}
