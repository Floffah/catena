import { IconPlus } from "@tabler/icons-react";
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import Link from "next/link";
import { notFound } from "next/navigation";

import PersonalAccessTokenList from "@/app/(site)/settings/tokens/PersonalAccessTokenList";
import { Button } from "@/components/ui/button";
import { $api } from "@/lib/api";
import { queryClient } from "@/lib/queryClient";
import { serverGetGitAccessTokens } from "@/lib/server/gitAccessTokens";

export default async function Page() {
    const tokens = await serverGetGitAccessTokens();

    if (!tokens) {
        return notFound();
    }

    const queryOptions = $api.queryOptions("get", "/v1/git-access-tokens");
    await queryClient.prefetchQuery({
        ...queryOptions,
        queryFn: () => tokens,
    });

    return (
        <HydrationBoundary state={dehydrate(queryClient)}>
            <main className="flex flex-1">
                <section className="flex w-full max-w-2xl flex-1 flex-col gap-4">
                    <header className="flex items-start justify-between gap-4">
                        <div className="flex flex-col gap-1">
                            <h2 className="text-xl font-bold">
                                Personal Access Tokens
                            </h2>
                            <p className="text-sm text-muted-foreground">
                                These tokens can currently only be used to
                                authenticate Git operations over HTTPS.
                            </p>
                        </div>

                        <Button asChild>
                            <Link href="/settings/tokens/new">
                                <IconPlus className="size-4" />
                                Create token
                            </Link>
                        </Button>
                    </header>

                    <PersonalAccessTokenList />
                </section>
            </main>
        </HydrationBoundary>
    );
}
