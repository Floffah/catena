import { IconArrowUp, IconFile, IconFolder } from "@tabler/icons-react";
import Link from "next/link";

import { serverGetRepositoryTree } from "@/lib/server/repository";

export async function RepositoryFileList({
    ownerName,
    repositoryName,
    branch,
    path,
}: {
    ownerName: string;
    repositoryName: string;
    branch?: string;
    path?: string;
}) {
    const tree = await serverGetRepositoryTree(
        ownerName,
        repositoryName,
        branch,
        path,
    );

    if (!tree) {
        return null;
    }

    const isRoot = tree.path === "";
    const parentPath = tree.path.split("/").slice(0, -1).join("/");
    const parentHref = parentPath
        ? `/${ownerName}/${repositoryName}/browse/${tree.ref}/${parentPath}`
        : `/${ownerName}/${repositoryName}/browse/${tree.ref}`;

    return (
        <div className="flex flex-col gap-4 rounded-lg bg-card text-card-foreground ring-1 ring-foreground/10">
            {(tree.entries.length > 0 || !isRoot) && (
                <ul className="divide-y divide-card-foreground/10 overflow-hidden rounded-lg">
                    {!isRoot && (
                        <li>
                            <Link
                                className="flex items-center gap-2 px-3 py-2 text-sm"
                                href={parentHref}
                            >
                                <IconArrowUp className="size-4 text-muted-foreground" />
                                <span>..</span>
                            </Link>
                        </li>
                    )}
                    {tree.entries.map((entry) => {
                        const Icon =
                            entry.type === "tree" ? IconFolder : IconFile;

                        return (
                            <li key={entry.path}>
                                <Link
                                    className="flex items-center gap-2 px-3 py-2 text-sm"
                                    href={`/${ownerName}/${repositoryName}/browse/${tree.ref}/${entry.path}`}
                                >
                                    <Icon className="size-4 text-muted-foreground" />
                                    <span>{entry.name}</span>
                                </Link>
                            </li>
                        );
                    })}
                </ul>
            )}
        </div>
    );
}
