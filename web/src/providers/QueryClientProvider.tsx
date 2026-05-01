"use client";

import { useAuth } from "@clerk/nextjs";
import {
    DefaultOptions,
    QueryClient,
    QueryClientProvider as TRQProvider,
} from "@tanstack/react-query";
import { PropsWithChildren, useEffect } from "react";

const queryClient = new QueryClient();

export default function QueryClientProvider({ children }: PropsWithChildren) {
    const auth = useAuth();

    useEffect(() => {
        if (auth.isLoaded) {
            queryClient.refetchQueries();
        } else {
            queryClient.cancelQueries();
        }
    }, [auth.isLoaded, auth.isSignedIn]);

    return <TRQProvider client={queryClient}>{children}</TRQProvider>;
}
