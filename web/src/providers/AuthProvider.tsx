"use client";

import { useAuth } from "@clerk/nextjs";
import { useQuery } from "@tanstack/react-query";
import { PropsWithChildren, createContext, useContext, useEffect } from "react";

import { setAuthToken } from "@/lib/api";

interface AuthContextValue {
    isLoading: boolean;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextValue>(null!);

export default function AuthProvider({ children }: PropsWithChildren) {
    const { getToken, isLoaded, isSignedIn, sessionId } = useAuth();

    const getTokenQuery = useQuery({
        queryKey: ["AuthProvider", "getToken", sessionId],
        queryFn: async ({ client }) => {
            if (!isLoaded || !isSignedIn) {
                setAuthToken(null);
                return null;
            }

            const newToken = await getToken();
            setAuthToken(newToken);

            client.refetchQueries({
                predicate: (query) => !!query.meta?.refetchOnAuth,
            });

            return newToken;
        },
        enabled: isLoaded && isSignedIn,
    });

    const isLoading = !isLoaded || (isSignedIn && getTokenQuery.isPending);
    const isAuthenticated = !!isSignedIn && !!getTokenQuery.data;

    useEffect(() => {
        if (isLoaded && !isSignedIn) {
            setAuthToken(null);
        }
    }, [isLoaded, isSignedIn]);

    return (
        <AuthContext.Provider value={{ isLoading, isAuthenticated }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useCatenaAuth() {
    const context = useContext(AuthContext);

    if (!context) {
        throw new Error("useCatenaAuth must be used within an AuthProvider");
    }

    return context;
}
