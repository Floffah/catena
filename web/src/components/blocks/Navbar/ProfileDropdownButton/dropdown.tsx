"use client";

import { SignOutButton, useClerk } from "@clerk/nextjs";
import { IconLogout, IconUserCog } from "@tabler/icons-react";
import Link from "next/link";

import UserAvatar from "@/components/UserAvatar";
import UserProfileDialogButton from "@/components/UserProfileDialogButton";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SchemaUser } from "@/types/api";

export default function ProfileDropdown({
    user: { name, avatarUrl },
}: {
    user: Partial<SchemaUser> & { name: string };
}) {
    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <button>
                    <UserAvatar user={{ name, avatarUrl }} />
                </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-full">
                <DropdownMenuItem asChild>
                    <Link
                        href={`/${name}`}
                        className="flex items-center gap-2 px-1 py-1.5 text-left text-sm"
                    >
                        <UserAvatar user={{ name, avatarUrl }} />
                        <div className="grid flex-1 text-left text-sm leading-tight">
                            <span className="truncate font-medium">{name}</span>
                            <span className="truncate text-xs">Profile</span>
                        </div>
                    </Link>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                    <DropdownMenuItem asChild>
                        <UserProfileDialogButton>
                            <IconUserCog className="size-4" />
                            Manage Account
                        </UserProfileDialogButton>
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
