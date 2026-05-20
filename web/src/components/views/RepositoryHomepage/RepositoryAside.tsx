import { IconClock } from "@tabler/icons-react";
import { formatRelative } from "date-fns";

import {
    serverGetRepository,
    serverGetRepositoryLatestCommit,
} from "@/lib/server/repository";

export default async function RepositoryAside({
    ownerName,
    repoName,
}: {
    ownerName: string;
    repoName: string;
}) {
    const [repo, latestCommit] = await Promise.all([
        serverGetRepository(ownerName, repoName),
        serverGetRepositoryLatestCommit(ownerName, repoName),
    ]);

    return (
        <aside className="flex max-w-80 shrink-0 flex-col gap-4">
            <section className="flex flex-col gap-1">
                <h2 className="text-lg font-semibold">About</h2>
                <p className="text-sm text-muted-foreground">
                    {repo?.description || (
                        <span className="italic">None provided</span>
                    )}
                </p>
            </section>

            <section className="flex flex-col gap-1">
                {latestCommit && (
                    <p className="flex items-center gap-1 text-sm text-muted-foreground">
                        <IconClock className="size-5" />
                        Last updated{" "}
                        {formatRelative(
                            new Date(latestCommit.authoredAt ?? ""),
                            new Date(),
                        )}
                    </p>
                )}
            </section>
        </aside>
    );
}
