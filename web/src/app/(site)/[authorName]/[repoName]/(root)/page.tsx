import { notFound } from "next/navigation";

import RepositoryHomepage from "@/components/views/RepositoryHomepage";
import { serverGetRepository } from "@/lib/server/repository";

export default async function Page({
    params,
}: {
    params: Promise<{ authorName: string; repoName: string }>;
}) {
    const { authorName, repoName } = await params;

    const repo = await serverGetRepository(authorName, repoName);

    if (!repo) {
        return notFound();
    }

    return (
        <RepositoryHomepage
            authorName={authorName}
            repoName={repoName}
            branch={repo.defaultBranch}
        />
    );
}
