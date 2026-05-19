import { Suspense } from "react";

import RepositoryAside from "./RepositoryAside";
import { RepositoryFileList } from "./RepositoryFileList";
import { RepositoryReadme } from "./RepositoryReadme";

export default async function RepositoryHomepage({
    ownerName,
    repoName,
    currentRef,
}: {
    ownerName: string;
    repoName: string;
    currentRef: string;
}) {
    return (
        <main className="flex gap-4">
            <div className="flex flex-1 flex-col gap-4">
                <RepositoryReadme
                    ownerName={ownerName}
                    repoName={repoName}
                    currentRef={currentRef}
                />
                <Suspense fallback={null}>
                    <RepositoryFileList
                        ownerName={ownerName}
                        repoName={repoName}
                        currentRef={currentRef}
                    />
                </Suspense>
            </div>
            <RepositoryAside ownerName={ownerName} repoName={repoName} />
        </main>
    );
}
