import { SignInButton } from "@clerk/nextjs";
import { auth } from "@clerk/nextjs/server";
import { IconPlus } from "@tabler/icons-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";

export default async function NewRepoButton() {
    const { isAuthenticated } = await auth();

    if (!isAuthenticated) {
        return (
            <SignInButton>
                <Button asChild variant="secondary">
                    <IconPlus />
                    New Repository
                </Button>
            </SignInButton>
        );
    }

    return (
        <Button asChild variant="secondary">
            <Link href="/new/repository">
                <IconPlus />
                New Repository
            </Link>
        </Button>
    );
}
