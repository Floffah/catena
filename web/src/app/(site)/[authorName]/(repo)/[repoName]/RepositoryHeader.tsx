import { IconBrandGit, IconGitBranch } from "@tabler/icons-react";

import { serverGetRepository } from "@/lib/server/repository";

export async function RepositoryHeader({
    ownerName,
    repositoryName,
}: {
    ownerName: string;
    repositoryName: string;
}) {
    const repo = await serverGetRepository(ownerName, repositoryName);

    return (
        <div className="flex items-center gap-4">
            <h1 className="flex items-center gap-1 text-xl">
                <IconBrandGit />
                {repo?.ownerName}
                <span className="text-2xl">/</span>
                {repo?.name}
            </h1>
        </div>
    );
}
