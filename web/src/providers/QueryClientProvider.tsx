"use client";

import {
    QueryClient,
    QueryClientProvider as TRQProvider,
} from "@tanstack/react-query";
import { PropsWithChildren } from "react";

const queryClient = new QueryClient();

export default function QueryClientProvider({ children }: PropsWithChildren) {
    return <TRQProvider client={queryClient}>{children}</TRQProvider>;
}
