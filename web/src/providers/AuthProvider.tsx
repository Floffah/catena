"use client";

import { useAuth } from "@clerk/nextjs";
import { useQueryClient } from "@tanstack/react-query";
import { Middleware } from "openapi-fetch";
import { PropsWithChildren, useEffect } from "react";

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

export default function AuthProvider({ children }: PropsWithChildren) {
    const auth = useAuth();
    const queryClient = useQueryClient();

    useEffect(() => {
        apiFetch.use(authMiddleware);

        return () => {
            apiFetch.eject(authMiddleware);
        };
    }, []);

    useEffect(() => {
        if (auth.isSignedIn) {
            auth.getToken().then((t) => {
                const needsInvalidate = !token && t;

                token = t;

                if (needsInvalidate) {
                    queryClient.clear();
                }
            });
        } else {
            token = null;
        }
    }, [auth, queryClient]);

    return children;
}
