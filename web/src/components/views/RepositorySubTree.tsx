import { Suspense } from "react";

import { RepositoryFileList } from "@/components/views/RepositoryHomepage/RepositoryFileList";
import { RepositoryReadme } from "@/components/views/RepositoryHomepage/RepositoryReadme";

export default async function RepositorySubTree({
    ownerName,
    repoName,
    branch,
    path,
}: {
    ownerName: string;
    repoName: string;
    branch: string;
    path: string;
}) {
    return (
        <main className="flex flex-col gap-4">
            <RepositoryFileList
                ownerName={ownerName}
                repoName={repoName}
                currentRef={branch}
                path={path}
            />

            <Suspense fallback={null}>
                <RepositoryReadme
                    ownerName={ownerName}
                    repoName={repoName}
                    currentRef={branch}
                    path={path}
                />
            </Suspense>
        </main>
    );
}
