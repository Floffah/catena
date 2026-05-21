import { notFound } from "next/navigation";

import RepositoryFileViewer from "@/components/views/RepositoryFileViewer";
import RepositoryHomepage from "@/components/views/RepositoryHomepage";
import RepositorySubTree from "@/components/views/RepositorySubTree";
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
        if (resolvedPath.pathType === "tree") {
            return (
                <RepositorySubTree
                    ownerName={ownerName}
                    repoName={repoName}
                    branch={resolvedPath.ref}
                    path={resolvedPath.path}
                />
            );
        }

        if (resolvedPath.pathType === "blob") {
            return (
                <RepositoryFileViewer
                    ownerName={ownerName}
                    repoName={repoName}
                    branch={resolvedPath.ref}
                    path={resolvedPath.path}
                />
            );
        }

        return notFound();
    }

    return (
        <RepositoryHomepage
            ownerName={ownerName}
            repoName={repoName}
            currentRef={resolvedPath.ref}
        />
    );
}
