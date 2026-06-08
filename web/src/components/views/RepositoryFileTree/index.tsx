import { preloadFileTree, prepareFileTreeInput } from "@pierre/trees";

import RepositoryFileTreeInner from "@/components/views/RepositoryFileTree/inner";
import { serverGetRepositoryTree } from "@/lib/server/repository";

export function RepositoryFileTreeFallback() {
    return (
        <aside
            aria-label="Repository files"
            className="h-72 overflow-hidden rounded-lg bg-card p-2 ring-1 ring-foreground/10 lg:sticky lg:top-4 lg:h-[calc(100vh-8rem)]"
        ></aside>
    );
}

export default async function RepositoryFileTree({
    ownerName,
    repoName,
    branch,
    path,
    isTree,
}: {
    ownerName: string;
    repoName: string;
    branch: string;
    path: string;
    isTree: boolean;
}) {
    console.log(path);
    const tree = await serverGetRepositoryTree(
        ownerName,
        repoName,
        branch,
        "/",
        true,
    );

    if (!tree) {
        return null;
    }

    const preparedInput = prepareFileTreeInput(
        tree.entries
            .filter((entry) => entry.type === "blob")
            .map((entry) => entry.path)
            .map(String),
    );

    const initialExpandedPaths = path
        .split("/")
        .slice(0, -1)
        .map((_, index, arr) => arr.slice(0, index + 1).join("/"));

    const lastSeg = path.split("/").pop();
    if (isTree && lastSeg) {
        initialExpandedPaths.push(lastSeg);
    }

    const payload = preloadFileTree({
        preparedInput,
        id: "project-tree-" + ownerName + "-" + repoName + "-" + branch,
        search: true,
        initialVisibleRowCount: 11,
        flattenEmptyDirectories: true,
        initialSelectedPaths: [path],
        initialExpandedPaths,
        initialExpansion: "open",
    });

    return (
        <RepositoryFileTreeInner
            preparedInput={preparedInput}
            preloadedData={payload}
            initialExpandedPaths={initialExpandedPaths}
            ownerName={ownerName}
            repoName={repoName}
            branch={branch}
            path={path}
        />
    );
}
