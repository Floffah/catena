import { SignInButton } from "@clerk/nextjs";
import { auth } from "@clerk/nextjs/server";
import { IconUser } from "@tabler/icons-react";
import { Suspense } from "react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Spinner } from "@/components/ui/spinner";
import { serverGetUserForClerkID } from "@/lib/server/users";

import ProfileDropdown from "./dropdown";

export default function ProfileDropdownButton() {
    return (
        <Suspense
            fallback={
                <Avatar>
                    <AvatarFallback>
                        <Spinner className="size-4" />
                    </AvatarFallback>
                </Avatar>
            }
        >
            <Inner />
        </Suspense>
    );
}

async function Inner() {
    const { isAuthenticated, userId } = await auth();

    const fallback = (
        <SignInButton>
            <Avatar asChild>
                <button>
                    <AvatarFallback>
                        <IconUser className="size-4" />
                    </AvatarFallback>
                </button>
            </Avatar>
        </SignInButton>
    );

    if (!userId) {
        return fallback;
    }

    const user = await serverGetUserForClerkID(userId);

    if (!isAuthenticated || !user) {
        return fallback;
    }

    return <ProfileDropdown userName={user.name} avatarUrl={user.avatarUrl} />;
}
