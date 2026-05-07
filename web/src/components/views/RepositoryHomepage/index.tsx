import { Suspense } from "react";

import RepositoryAside from "./RepositoryAside";
import { RepositoryFileList } from "./RepositoryFileList";
import { RepositoryReadme } from "./RepositoryReadme";

export default async function RepositoryHomepage({
    authorName,
    repoName,
    branch,
}: {
    authorName: string;
    repoName: string;
    branch: string;
}) {
    return (
        <main className="flex gap-4">
            <div className="flex flex-1 flex-col gap-4">
                <RepositoryReadme
                    ownerName={authorName}
                    repositoryName={repoName}
                    branch={branch}
                />
                <Suspense fallback={null}>
                    <RepositoryFileList
                        ownerName={authorName}
                        repositoryName={repoName}
                        branch={branch}
                    />
                </Suspense>
            </div>
            <RepositoryAside ownerName={authorName} repositoryName={repoName} />
        </main>
    );
}
