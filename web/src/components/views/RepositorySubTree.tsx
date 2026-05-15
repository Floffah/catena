import { Suspense } from "react";

import { RepositoryFileList } from "@/components/views/RepositoryHomepage/RepositoryFileList";
import { RepositoryReadme } from "@/components/views/RepositoryHomepage/RepositoryReadme";

export default async function RepositorySubTree({
    authorName,
    repoName,
    branch,
    path,
}: {
    authorName: string;
    repoName: string;
    branch: string;
    path: string;
}) {
    return (
        <main className="flex flex-col gap-4">
            <Suspense fallback={null}>
                <RepositoryFileList
                    ownerName={authorName}
                    repositoryName={repoName}
                    branch={branch}
                    path={path}
                />
            </Suspense>

            <Suspense fallback={null}>
                <RepositoryReadme
                    ownerName={authorName}
                    repositoryName={repoName}
                    branch={branch}
                    path={path}
                />
            </Suspense>
        </main>
    );
}
