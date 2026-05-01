"use client";

import { $api } from "@/lib/api";

export default function Page() {
    const health = $api.useQuery("get", "/v1/healthz");

    return null;
}
