import { IconCornerDownRight } from "@tabler/icons-react";
import { notFound } from "next/navigation";
import { PropsWithChildren } from "react";

import { RepositoryHeader } from "@/components/views/RepositoryHomepage/RepositoryHeader";
import {
    serverGetRepository,
    serverResolveRepositoryGitPath,
} from "@/lib/server/repository";

export default async function RepoLayout({
    children,
    authorName,
    repoName,
    path,
}: PropsWithChildren<{
    authorName: string;
    repoName: string;
    path?: string[];
}>) {
    const repo = await serverGetRepository(authorName, repoName);

    if (!repo) {
        return notFound();
    }

    let currentRef = repo.defaultBranch;
    console.log(path);
    if (path) {
        const resolvedPath = await serverResolveRepositoryGitPath(
            authorName,
            repoName,
            path.join("/"),
        );

        if (resolvedPath) {
            currentRef = resolvedPath.ref;
        } else {
            return notFound();
        }
    }

    return (
        <div className="container mx-auto flex flex-1 flex-col gap-4 p-4">
            <div className="flex flex-col gap-1">
                <RepositoryHeader
                    ownerName={authorName}
                    repositoryName={repoName}
                    currentRef={currentRef}
                />
                <p className="flex items-center gap-1 text-sm text-muted-foreground underline">
                    <IconCornerDownRight className="size-4" />
                    <a href="#browse">
                        Scroll down to browse the repository contents
                    </a>
                </p>
            </div>

            {children}
        </div>
    );
}
