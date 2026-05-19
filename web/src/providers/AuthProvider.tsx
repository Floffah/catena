"use client";

import { useAuth } from "@clerk/nextjs";
import { Middleware } from "openapi-fetch";
import { PropsWithChildren, createContext, useContext, useEffect } from "react";

import { apiFetch } from "@/lib/api";

interface AuthContextValue {
    isLoading: boolean;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextValue>(null!);

export default function AuthProvider({ children }: PropsWithChildren) {
    const { getToken, isLoaded, isSignedIn } = useAuth();

    useEffect(() => {
        const authMiddleware: Middleware = {
            async onRequest({ request }) {
                const token = await getToken();

                if (token) {
                    request.headers.set("Authorization", `Bearer ${token}`);
                }

                return request;
            },
        };

        apiFetch.use(authMiddleware);

        return () => {
            apiFetch.eject(authMiddleware);
        };
    }, [getToken]);

    return (
        <AuthContext.Provider
            value={{ isLoading: !isLoaded, isAuthenticated: !!isSignedIn }}
        >
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
