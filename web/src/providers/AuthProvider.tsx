"use client";

import { useAuth } from "@clerk/nextjs";
import { useQuery } from "@tanstack/react-query";
import { Middleware } from "openapi-fetch";
import { PropsWithChildren, createContext, useContext, useEffect } from "react";

import { apiFetch } from "@/lib/api";

let token: string | null = null;

const authMiddleware: Middleware = {
    async onRequest({ request }) {
        if (token) {
            request.headers.set("Authorization", `Bearer ${token}`);
        }

        return request;
    },
};

interface AuthContextValue {
    isLoading: boolean;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextValue>(null!);

export default function AuthProvider({ children }: PropsWithChildren) {
    const { getToken, isLoaded, isSignedIn, sessionId } = useAuth();

    const getTokenQuery = useQuery({
        queryKey: ["AuthProvider", "getToken", sessionId],
        queryFn: async () => {
            if (!isLoaded || !isSignedIn) {
                token = null;
                return null;
            }

            const newToken = await getToken();
            token = newToken;
            return newToken;
        },
        enabled: isLoaded && isSignedIn,
    });

    const isLoading = !isLoaded || (isSignedIn && getTokenQuery.isPending);
    const isAuthenticated = !!isSignedIn && !!getTokenQuery.data;

    useEffect(() => {
        if (isLoaded && !isSignedIn) {
            token = null;
        }
    }, [isLoaded, isSignedIn]);

    useEffect(() => {
        apiFetch.use(authMiddleware);

        return () => {
            apiFetch.eject(authMiddleware);
        };
    }, []);

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
