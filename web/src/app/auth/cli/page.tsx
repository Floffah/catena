"use client";

import { useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { $api } from "@/lib/api";

export default function Page() {
    const params = useSearchParams();
    const redirectUrl = params.get("redirect_uri");

    const getTokenMutation = $api.useMutation(
        "post",
        "/v1/auth/clerk/sign-in-token",
        {
            onSuccess: ({ token }) => {
                if (redirectUrl) {
                    const url = new URL(redirectUrl);
                    url.searchParams.set("token", token);
                    window.location.href = url.toString();
                }
            },
        },
    );

    if (!redirectUrl) {
        return (
            <main className="flex flex-1 flex-col items-center justify-center gap-1">
                <p>
                    No redirect URL provided. Please provide a redirect URL to
                    authenticate Catena CLI.
                </p>
            </main>
        );
    }

    return (
        <main className="flex flex-1 flex-col items-center justify-center gap-1">
            <p>
                Are you sure you want to authenticate{" "}
                <span className="font-bold">Catena CLI</span> with your account?
            </p>
            <Button onClick={() => getTokenMutation.mutate({})}>
                Yes, authenticate Catena CLI
            </Button>
        </main>
    );
}
