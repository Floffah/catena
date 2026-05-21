import { notFound } from "next/navigation";

import RepositorySettingsForm from "@/components/views/RepositorySettingsForm";
import {
    serverGetRepository,
    serverListRepositoryRefs,
} from "@/lib/server/repository";
import { serverGetCurrentUser } from "@/lib/server/users";

export default async function Page({
    params,
}: {
    params: Promise<{ ownerName: string; repoName: string }>;
}) {
    const { ownerName, repoName } = await params;
    const [repository, user] = await Promise.all([
        serverGetRepository(ownerName, repoName),
        serverGetCurrentUser(),
    ]);

    if (!repository || !user || user.name !== repository.ownerName) {
        return notFound();
    }

    const refs = await serverListRepositoryRefs(ownerName, repoName);
    const branchNames = Array.from(
        new Set([
            repository.defaultBranch,
            ...(refs?.refs.map((ref) => ref.name) ?? []),
        ]),
    );

    return (
        <RepositorySettingsForm
            branchNames={branchNames}
            repository={repository}
        />
    );
}
