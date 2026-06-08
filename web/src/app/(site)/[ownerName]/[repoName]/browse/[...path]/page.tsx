import { notFound } from "next/navigation";

import RepositoryBrowseView from "@/components/views/RepositoryBrowseView";
import RepositoryHomepage from "@/components/views/RepositoryHomepage";
import { serverResolveRepositoryGitPath } from "@/lib/server/repository";

export default async function Page({
    params,
}: {
    params: Promise<{
        ownerName: string;
        repoName: string;
        path: string[];
    }>;
}) {
    const { ownerName, repoName, path } = await params;
    const resolvedPath = await serverResolveRepositoryGitPath(
        ownerName,
        repoName,
        path.join("/"),
    );

    if (!resolvedPath) {
        return notFound();
    }

    if (resolvedPath.pathType !== "root") {
        return (
            <RepositoryBrowseView
                ownerName={ownerName}
                repoName={repoName}
                branch={resolvedPath.ref}
                path={resolvedPath.path}
                pathType={resolvedPath.pathType}
            />
        );
    }

    return (
        <RepositoryHomepage
            ownerName={ownerName}
            repoName={repoName}
            currentRef={resolvedPath.ref}
        />
    );
}
