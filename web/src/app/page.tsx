"use client";

import { SignInButton } from "@clerk/nextjs";

import { Button } from "@/components/ui/button";

export default function Home() {
    return (
        <div className="flex flex-1 items-center justify-center">
            <SignInButton>
                <Button>Sign in</Button>
            </SignInButton>
        </div>
    );
}
