import { notFound } from "next/navigation";

import RepositoryHomepage from "@/components/views/RepositoryHomepage";
import RepositorySubTree from "@/components/views/RepositorySubTree";
import { serverResolveRepositoryGitPath } from "@/lib/server/repository";

export default async function Page({
    params,
}: {
    params: Promise<{
        authorName: string;
        repoName: string;
        path: string[];
    }>;
}) {
    const { authorName, repoName, path } = await params;
    const resolvedPath = await serverResolveRepositoryGitPath(
        authorName,
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
                    authorName={authorName}
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
            authorName={authorName}
            repoName={repoName}
            branch={resolvedPath.ref}
        />
    );
}
