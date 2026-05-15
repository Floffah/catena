"use client";

import { QueryClientProvider as TRQProvider } from "@tanstack/react-query";
import { PropsWithChildren } from "react";

import { createQueryClient } from "@/lib/queryClient";

const queryClient = createQueryClient();

export default function QueryClientProvider({ children }: PropsWithChildren) {
    return <TRQProvider client={queryClient}>{children}</TRQProvider>;
}
