import {
    IconBook,
    IconGitBranch,
    IconLock,
    IconWorld,
} from "@tabler/icons-react";
import Link from "next/link";

import {
    Empty,
    EmptyDescription,
    EmptyHeader,
    EmptyMedia,
    EmptyTitle,
} from "@/components/ui/empty";
import { serverListFeaturedRepositoriesForUser } from "@/lib/server/users";

export function ProfileRepositoriesSkeleton() {
    return (
        <section className="flex flex-col gap-4" aria-label="Repositories">
            <div className="flex items-end justify-between gap-4">
                <div className="space-y-2">
                    <div className="h-6 w-32 animate-pulse rounded-md bg-muted" />
                    <div className="h-4 w-56 animate-pulse rounded-md bg-muted" />
                </div>
            </div>

            <div className="grid gap-3 md:grid-cols-2">
                {Array.from({ length: 4 }).map((_, index) => (
                    <div
                        key={index}
                        className="min-h-36 animate-pulse rounded-lg border border-border bg-muted/30"
                    />
                ))}
            </div>
        </section>
    );
}

export default async function ProfileRepositories({
    ownerName,
}: {
    ownerName: string;
}) {
    const repositories = await serverListFeaturedRepositoriesForUser(ownerName);

    if (!repositories) {
        return null;
    }

    return (
        <section className="flex flex-col gap-4" aria-labelledby="repositories">
            <header className="flex flex-col gap-1">
                <h2
                    id="repositories"
                    className="font-heading text-2xl font-bold tracking-tight"
                >
                    Repositories
                </h2>
                <p className="text-sm text-muted-foreground">
                    Featured work from @{ownerName}.
                </p>
            </header>

            {repositories.repositories.length === 0 ? (
                <Empty className="border border-border bg-muted/20">
                    <EmptyHeader>
                        <EmptyMedia variant="icon">
                            <IconBook />
                        </EmptyMedia>
                        <EmptyTitle>No repositories yet</EmptyTitle>
                        <EmptyDescription>
                            Repositories owned by this user will show up here.
                        </EmptyDescription>
                    </EmptyHeader>
                </Empty>
            ) : (
                <div className="grid gap-3 md:grid-cols-2">
                    {repositories.repositories.map((repository) => (
                        <Link
                            key={repository.id}
                            href={`/${repository.ownerName}/${repository.name}`}
                            className="group min-w-0 rounded-lg border border-border bg-background p-4 transition-colors hover:border-primary/50 hover:bg-muted/35 focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
                        >
                            <article className="flex flex-col gap-4">
                                <div className="flex min-w-0 items-start justify-between gap-3">
                                    <div className="min-w-0">
                                        <h3 className="truncate font-heading text-base font-semibold tracking-tight text-foreground group-hover:text-primary">
                                            {repository.name}
                                        </h3>
                                        <p className="mt-1 flex items-center gap-1.5 text-xs text-muted-foreground">
                                            {repository.visibility ===
                                            "private" ? (
                                                <IconLock className="size-3.5" />
                                            ) : (
                                                <IconWorld className="size-3.5" />
                                            )}
                                            <span className="capitalize">
                                                {repository.visibility}
                                            </span>
                                        </p>
                                    </div>

                                    <IconGitBranch className="mt-1 size-4 shrink-0 text-muted-foreground" />
                                </div>

                                {repository.description && (
                                    <p className="line-clamp-2 text-sm/relaxed text-muted-foreground">
                                        {repository.description}
                                    </p>
                                )}
                            </article>
                        </Link>
                    ))}
                </div>
            )}
        </section>
    );
}
