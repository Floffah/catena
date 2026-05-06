import { IconCornerDownRight } from "@tabler/icons-react";

import RepositoryAside from "@/app/(site)/[authorName]/(repo)/[repoName]/RepositoryAside";
import { RepositoryHeader } from "@/app/(site)/[authorName]/(repo)/[repoName]/RepositoryHeader";
import { RepositoryReadme } from "@/app/(site)/[authorName]/(repo)/[repoName]/RepositoryReadme";

export default async function Page({
    params,
}: {
    params: Promise<{ authorName: string; repoName: string }>;
}) {
    const { authorName, repoName } = await params;

    return (
        <div className="container mx-auto flex flex-1 flex-col gap-4 p-4">
            <div className="flex flex-col gap-1">
                <RepositoryHeader
                    ownerName={authorName}
                    repositoryName={repoName}
                />
                <p className="flex items-center gap-1 text-sm text-muted-foreground underline">
                    <IconCornerDownRight className="size-4" />
                    <a href="#browse">
                        Scroll down to browse the repository contents
                    </a>
                </p>
            </div>

            <main className="flex gap-4">
                <div className="flex flex-1 flex-col gap-4">
                    <RepositoryReadme
                        ownerName={authorName}
                        repositoryName={repoName}
                    />
                </div>
                <RepositoryAside
                    ownerName={authorName}
                    repositoryName={repoName}
                />
            </main>
        </div>
    );
}
