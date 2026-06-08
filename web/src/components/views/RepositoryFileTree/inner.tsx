"use client";

import { FileTreePreparedInput, FileTreeSsrPayload } from "@pierre/trees";
import { FileTree, useFileTree } from "@pierre/trees/react";
import { useRouter } from "next/navigation";

export default function RepositoryFileTreeInner({
    preparedInput,
    preloadedData,
    initialExpandedPaths,
    ownerName,
    repoName,
    branch,
    path,
}: {
    preparedInput: FileTreePreparedInput;
    preloadedData: FileTreeSsrPayload;
    initialExpandedPaths: string[];
    ownerName: string;
    repoName: string;
    branch: string;
    path: string;
}) {
    const router = useRouter();
    const { model } = useFileTree({
        preparedInput,
        id: preloadedData.id,
        search: true,
        initialVisibleRowCount: 11,
        flattenEmptyDirectories: true,
        initialSelectedPaths: [path],
        initialExpandedPaths,
        initialExpansion: "open",
        onSelectionChange: (paths) => {
            if (paths.length === 0) {
                return;
            }

            const selectedPath = paths[0];
            router.push(
                `/${ownerName}/${repoName}/browse/${branch}/${selectedPath}`,
            );
        },
    });

    return (
        <aside
            aria-label="Repository files"
            className="overflow-y-auto rounded-lg bg-card p-2 ring-1 ring-foreground/10 lg:sticky lg:top-4"
        >
            <FileTree model={model} preloadedData={preloadedData} />
        </aside>
    );
}
