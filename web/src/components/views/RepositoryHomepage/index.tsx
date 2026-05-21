import { Suspense } from "react";

import { RepositoryHeader } from "@/components/views/RepositoryHomepage/RepositoryHeader";

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
        <main className="flex flex-col gap-4 lg:flex-row">
            <div className="flex flex-1 flex-col gap-4">
                <RepositoryHeader
                    ownerName={ownerName}
                    repoName={repoName}
                    currentRef={currentRef}
                />
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
