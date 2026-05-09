"use client";

import {
    SignInButton,
    SignOutButton,
    UserButton,
    UserProfile,
    useClerk,
} from "@clerk/nextjs";
import { auth } from "@clerk/nextjs/server";
import { IconLogout, IconUser, IconUserCog } from "@tabler/icons-react";
import Link from "next/link";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export default function ProfileDropdown({
    userName,
    avatarUrl,
}: {
    userName: string;
    avatarUrl?: string | null;
}) {
    const { openUserProfile } = useClerk();

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Avatar asChild>
                    <button>
                        <AvatarFallback>
                            {userName[0].toUpperCase() ?? (
                                <IconUser className="size-4" />
                            )}
                        </AvatarFallback>
                        {avatarUrl && (
                            <AvatarImage src={avatarUrl} alt={userName} />
                        )}
                    </button>
                </Avatar>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
                <DropdownMenuItem asChild>
                    <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                        <Avatar className="h-8 w-8 rounded-lg">
                            {avatarUrl && (
                                <AvatarImage src={avatarUrl} alt={userName} />
                            )}
                            <AvatarFallback className="rounded-lg">
                                {userName[0].toUpperCase() ?? (
                                    <IconUser className="size-4" />
                                )}
                            </AvatarFallback>
                        </Avatar>
                        <div className="grid flex-1 text-left text-sm leading-tight">
                            <span className="truncate font-medium">
                                {userName}
                            </span>
                            <span className="truncate text-xs">Profile</span>
                        </div>
                    </div>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                    <DropdownMenuItem asChild>
                        <button onClick={() => openUserProfile()}>
                            <IconUserCog className="size-4" />
                            Manage Account
                        </button>
                    </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                    <SignOutButton>
                        <DropdownMenuItem>
                            <IconLogout className="size-4" />
                            Sign Out
                        </DropdownMenuItem>
                    </SignOutButton>
                </DropdownMenuGroup>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
