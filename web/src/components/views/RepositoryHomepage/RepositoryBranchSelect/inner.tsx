"use client";

import { IconSelector } from "@tabler/icons-react";
import Link from "next/link";
import { useParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export default function RepositoryBranchSelectInner({
    ownerName,
    repoName,
    currentRef,
    availableBranches,
}: {
    ownerName: string;
    repoName: string;
    currentRef: string;
    availableBranches?: string[];
}) {
    const params = useParams();

    const browsePathSegments = params["path"] as string[];
    const browsePath = browsePathSegments ? browsePathSegments.join("/") : "";

    if (!availableBranches) {
        return currentRef;
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="outline">
                    <span className="text-muted-foreground">Currently on</span>{" "}
                    <strong>{currentRef}</strong>
                    <IconSelector className="size-4" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
                <DropdownMenuGroup>
                    {availableBranches.map((branch) => {
                        const pathBase = `/${ownerName}/${repoName}/browse`;
                        let href = `${pathBase}/${branch}`;

                        if (browsePath) {
                            href = `${pathBase}/${browsePath.replace(currentRef, branch)}`;
                        }

                        return (
                            <DropdownMenuItem asChild key={branch}>
                                <Link href={href}>{branch}</Link>
                            </DropdownMenuItem>
                        );
                    })}
                </DropdownMenuGroup>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
