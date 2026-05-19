import { IconCircleDot, IconPlus } from "@tabler/icons-react";
import Link from "next/link";
import { notFound } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
    Empty,
    EmptyDescription,
    EmptyHeader,
    EmptyMedia,
    EmptyTitle,
} from "@/components/ui/empty";
import { serverListRepositoryIssues } from "@/lib/server/repository";

const statusLabel = (status: string) => status.replaceAll("_", " ");

export default async function Page({
    params,
}: {
    params: Promise<{ ownerName: string; repoName: string }>;
}) {
    const { ownerName, repoName } = await params;
    const issues = await serverListRepositoryIssues(ownerName, repoName);

    if (!issues) {
        return notFound();
    }

    return (
        <main className="flex flex-col gap-4">
            <header className="flex items-start justify-between gap-4">
                <div>
                    <h1 className="text-2xl font-bold">Issues</h1>
                    <p className="text-sm text-muted-foreground">
                        Track work, bugs, and ideas for this repository.
                    </p>
                </div>

                <Button asChild>
                    <Link href={`/${ownerName}/${repoName}/issues/new`}>
                        <IconPlus />
                        New issue
                    </Link>
                </Button>
            </header>

            {issues.issues.length === 0 ? (
                <Empty className="rounded-xl border border-border bg-card">
                    <EmptyHeader>
                        <EmptyMedia variant="icon">
                            <IconCircleDot />
                        </EmptyMedia>
                        <EmptyTitle>No issues yet</EmptyTitle>
                        <EmptyDescription>
                            Issues created for this repository will show up
                            here.
                        </EmptyDescription>
                    </EmptyHeader>
                </Empty>
            ) : (
                <section className="overflow-hidden rounded-xl border border-border bg-card">
                    {issues.issues.map((issue) => (
                        <Link
                            key={issue.id}
                            href={`/${ownerName}/${repoName}/issues/${issue.number}`}
                        >
                            <article className="flex items-start gap-3 border-b border-border px-4 py-3 last:border-b-0">
                                <IconCircleDot className="mt-0.5 size-4 text-muted-foreground" />
                                <div className="min-w-0 flex-1">
                                    <div className="flex flex-wrap items-baseline gap-x-2 gap-y-1">
                                        <h2 className="truncate text-sm font-medium">
                                            {issue.title}
                                        </h2>
                                        <span className="font-mono text-xs text-muted-foreground">
                                            {issue.reference}
                                        </span>
                                    </div>
                                    <p className="mt-1 text-xs text-muted-foreground capitalize">
                                        {statusLabel(issue.status)}
                                    </p>
                                </div>
                            </article>
                        </Link>
                    ))}
                </section>
            )}
        </main>
    );
}
