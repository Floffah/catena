import { notFound } from "next/navigation";
import { Suspense } from "react";

import RepositoryFileTree, {
    RepositoryFileTreeFallback,
} from "@/components/views/RepositoryFileTree";
import RepositoryFileViewer from "@/components/views/RepositoryFileViewer";
import RepositorySubTree from "@/components/views/RepositorySubTree";

export default function RepositoryBrowseView({
    ownerName,
    repoName,
    branch,
    path,
    pathType,
}: {
    ownerName: string;
    repoName: string;
    branch: string;
    path: string;
    pathType: "blob" | "commit" | "tree";
}) {
    if (pathType === "commit") {
        return notFound();
    }

    return (
        <div className="grid min-w-0 flex-1 gap-4 lg:grid-cols-[16rem_minmax(0,1fr)]">
            <Suspense fallback={<RepositoryFileTreeFallback />}>
                <RepositoryFileTree
                    ownerName={ownerName}
                    repoName={repoName}
                    branch={branch}
                    path={path}
                    isTree={pathType === "tree"}
                />
            </Suspense>
            <div className="min-w-0">
                {pathType === "tree" && (
                    <RepositorySubTree
                        ownerName={ownerName}
                        repoName={repoName}
                        branch={branch}
                        path={path}
                    />
                )}

                {pathType === "blob" && (
                    <RepositoryFileViewer
                        ownerName={ownerName}
                        repoName={repoName}
                        branch={branch}
                        path={path}
                    />
                )}
            </div>
        </div>
    );
}
