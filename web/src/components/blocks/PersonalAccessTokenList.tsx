"use client";

import { IconKey, IconTrash } from "@tabler/icons-react";

import { Button } from "@/components/ui/button";
import {
    Empty,
    EmptyDescription,
    EmptyHeader,
    EmptyMedia,
    EmptyTitle,
} from "@/components/ui/empty";
import { $api } from "@/lib/api";

const dateFormatter = new Intl.DateTimeFormat("en-GB", {
    dateStyle: "medium",
});

export default function PersonalAccessTokenList() {
    const { data: tokens = [], refetch: refetchTokens } = $api.useQuery(
        "get",
        "/v1/git-access-tokens",
    );

    const revokeTokenMutation = $api.useMutation(
        "delete",
        "/v1/git-access-tokens/{id}",
        {
            onSuccess: async () => {
                await refetchTokens();
            },
        },
    );

    if (tokens.length === 0) {
        return (
            <Empty className="rounded-lg border border-border bg-card">
                <EmptyHeader>
                    <EmptyMedia variant="icon">
                        <IconKey />
                    </EmptyMedia>
                    <EmptyTitle>No active personal access tokens</EmptyTitle>
                    <EmptyDescription>
                        Create a token when you need to authenticate Git
                        operations over HTTPS.
                    </EmptyDescription>
                </EmptyHeader>
            </Empty>
        );
    }

    return (
        <div className="divide-y divide-border overflow-hidden rounded-lg ring-1 ring-foreground/10">
            {tokens.map((token) => (
                <div
                    key={token.id}
                    className="flex items-center justify-between gap-4 px-4 py-3"
                >
                    <div className="flex min-w-0 flex-col gap-0.5">
                        <p className="truncate text-sm font-medium">
                            {token.name}
                        </p>
                        <p className="flex items-center gap-2 text-xs text-muted-foreground">
                            <span>Starts with {token.tokenPrefix}...</span>
                            <span>&bull;</span>
                            <span>
                                {token.expiresAt
                                    ? `Expires ${dateFormatter.format(
                                          new Date(token.expiresAt),
                                      )}`
                                    : "Never expires"}
                            </span>
                        </p>
                    </div>

                    <Button
                        aria-label={`Revoke ${token.name}`}
                        disabled={
                            revokeTokenMutation.isPending &&
                            revokeTokenMutation.variables?.params.path.id ===
                                token.id
                        }
                        onClick={() =>
                            revokeTokenMutation.mutate({
                                params: {
                                    path: {
                                        id: token.id,
                                    },
                                },
                            })
                        }
                        size="icon"
                        title="Revoke token"
                        variant="ghost"
                    >
                        <IconTrash className="size-4" />
                    </Button>
                </div>
            ))}
        </div>
    );
}
