"use client";

import { QueryClientProvider as TRQProvider } from "@tanstack/react-query";
import { PropsWithChildren } from "react";

import { queryClient } from "@/lib/queryClient";

export default function QueryClientProvider({ children }: PropsWithChildren) {
    return <TRQProvider client={queryClient}>{children}</TRQProvider>;
}
