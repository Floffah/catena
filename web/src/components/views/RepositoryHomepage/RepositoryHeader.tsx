import { IconBrandGit, IconError404 } from "@tabler/icons-react";
import { Suspense } from "react";

import RepositoryBranchSelect from "@/components/views/RepositoryHomepage/RepositoryBranchSelect";
import { serverGetRepository } from "@/lib/server/repository";

export async function RepositoryHeader({
    ownerName,
    repositoryName,
    currentRef,
}: {
    ownerName: string;
    repositoryName: string;
    currentRef: string;
}) {
    const repo = await serverGetRepository(ownerName, repositoryName);

    if (!repo) {
        return (
            <div className="flex items-center gap-4">
                <p className="flex items-center gap-1 text-xl">
                    <IconError404 className="size-5" />
                    No Repo
                </p>
            </div>
        );
    }

    return (
        <div className="flex items-center gap-4">
            <h1 className="flex items-center gap-1 text-xl">
                <IconBrandGit />
                {repo.ownerName}
                <span className="text-2xl">/</span>
                {repo.name}
                <span className="text-2xl">/</span>
                <Suspense fallback={currentRef}>
                    <RepositoryBranchSelect
                        ownerName={ownerName}
                        repositoryName={repositoryName}
                        currentRef={currentRef}
                    />
                </Suspense>
            </h1>
        </div>
    );
}
