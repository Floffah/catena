import { SignInButton } from "@clerk/nextjs";
import { auth } from "@clerk/nextjs/server";
import { Suspense } from "react";

import { UserAvatarFallback } from "@/components/UserAvatar";
import {
    serverGetCurrentUser,
    serverGetUserForClerkID,
} from "@/lib/server/users";

import ProfileDropdown from "./dropdown";

export default function ProfileDropdownButton() {
    return (
        <Suspense fallback={<UserAvatarFallback loading />}>
            <Inner />
        </Suspense>
    );
}

async function Inner() {
    // const { isAuthenticated, userId } = await auth();

    const fallback = (
        <SignInButton>
            <UserAvatarFallback />
        </SignInButton>
    );

    // if (!userId) {
    //     return fallback;
    // }

    const user = await serverGetCurrentUser();

    if (!user) {
        return fallback;
    }

    return <ProfileDropdown user={user} />;
}
