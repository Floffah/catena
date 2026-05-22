import { render, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, mock, test } from "bun:test";
import { HttpResponse, http } from "msw";
import { PropsWithChildren, useEffect } from "react";

import { server } from "@/test/server";

process.env.NEXT_PUBLIC_CATENA_INSTANCE_URL = "http://catena.test";

let token = "initial-token";

mock.module("@clerk/nextjs", () => ({
    useAuth() {
        return {
            getToken: async () => token,
            isLoaded: true,
            isSignedIn: true,
        };
    },
}));

const authHeaders: (string | null)[] = [];

afterEach(() => {
    authHeaders.length = 0;
    token = "initial-token";
});

describe("AuthProvider", () => {
    test("adds the latest Clerk token to API requests", async () => {
        server.use(
            http.get("http://catena.test/v1/user", ({ request }) => {
                authHeaders.push(request.headers.get("authorization"));

                return HttpResponse.json({
                    avatarUrl: null,
                    createdAt: "2026-05-22T00:00:00Z",
                    displayName: "Floffah",
                    id: "019deb10-dafc-743f-8cfc-289a80c13af1",
                    name: "floffah",
                    updatedAt: "2026-05-22T00:00:00Z",
                });
            }),
        );

        const authProviderModule = await import("./AuthProvider");
        const apiModule = await import("@/lib/api");
        const AuthProvider = authProviderModule.default;
        const { apiFetch } = apiModule;

        function ApiProbe({ requestKey }: { requestKey: string }) {
            useEffect(() => {
                void apiFetch.GET("/v1/user");
            }, [requestKey]);

            return null;
        }

        function Wrapper({ children }: PropsWithChildren) {
            return <AuthProvider>{children}</AuthProvider>;
        }

        const view = render(<ApiProbe requestKey="first" />, {
            wrapper: Wrapper,
        });

        await waitFor(() => {
            expect(authHeaders).toEqual(["Bearer initial-token"]);
        });

        token = "updated-token";
        view.rerender(<ApiProbe requestKey="second" />);

        await waitFor(() => {
            expect(authHeaders).toEqual([
                "Bearer initial-token",
                "Bearer updated-token",
            ]);
        });
    });
});
