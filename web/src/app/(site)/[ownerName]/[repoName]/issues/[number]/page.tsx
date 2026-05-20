import { IconArrowLeft, IconCircleDot } from "@tabler/icons-react";
import Link from "next/link";
import { notFound } from "next/navigation";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { serverGetRepositoryIssue } from "@/lib/server/repository";

const dateFormatter = new Intl.DateTimeFormat("en-GB", {
    dateStyle: "medium",
    timeStyle: "short",
});

const statusLabel = (status: string) => status.replaceAll("_", " ");

export default async function Page({
    params,
}: {
    params: Promise<{ ownerName: string; repoName: string; number: string }>;
}) {
    const { ownerName, repoName, number } = await params;
    const issueNumber = Number(number);

    if (!Number.isSafeInteger(issueNumber) || issueNumber < 1) {
        return notFound();
    }

    const issue = await serverGetRepositoryIssue(
        ownerName,
        repoName,
        issueNumber,
    );

    if (!issue) {
        return notFound();
    }

    return (
        <main className="mx-auto flex w-full max-w-4xl flex-col gap-6">
            <div>
                <Button asChild size="sm" variant="ghost">
                    <Link href={`/${ownerName}/${repoName}/issues`}>
                        <IconArrowLeft className="size-3" />
                        Back to issues
                    </Link>
                </Button>
            </div>

            <header className="flex flex-col gap-3">
                <div className="flex flex-wrap items-center gap-2">
                    <span className="inline-flex items-center gap-1 rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground capitalize">
                        <IconCircleDot className="size-3" />
                        {statusLabel(issue.status)}
                    </span>
                </div>

                <h1 className="flex items-center gap-4 text-3xl font-bold tracking-tight">
                    {issue.title}
                    <span className="rounded bg-muted px-2 py-1 font-mono text-sm text-muted-foreground">
                        {issue.reference}
                    </span>
                </h1>

                <p className="text-sm text-muted-foreground">
                    Opened {dateFormatter.format(new Date(issue.createdAt))}
                </p>
            </header>

            <Separator />

            {issue.body ? (
                <article className="prose prose-sm max-w-none dark:prose-invert">
                    <Markdown remarkPlugins={[remarkGfm]}>
                        {issue.body}
                    </Markdown>
                </article>
            ) : (
                <p className="text-sm text-muted-foreground">
                    No description provided.
                </p>
            )}
        </main>
    );
}
