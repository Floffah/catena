import { IconFolderOff, IconFolderPlus } from "@tabler/icons-react";
import { notFound } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
    Empty,
    EmptyContent,
    EmptyDescription,
    EmptyHeader,
    EmptyMedia,
    EmptyTitle,
} from "@/components/ui/empty";
import RepositoryHomepage from "@/components/views/RepositoryHomepage";
import {
    serverGetRepository,
    serverListRepositoryRefs,
} from "@/lib/server/repository";

export default async function Page({
    params,
}: {
    params: Promise<{ ownerName: string; repoName: string }>;
}) {
    const { ownerName, repoName } = await params;

    const repo = await serverGetRepository(ownerName, repoName);

    if (!repo) {
        return notFound();
    }

    const refs = await serverListRepositoryRefs(ownerName, repoName);

    if (!refs || refs.refs.length === 0) {
        return (
            <Empty className="mx-auto mt-10">
                <EmptyHeader>
                    <EmptyMedia variant="icon">
                        <IconFolderOff />
                    </EmptyMedia>
                    <EmptyTitle>Repository empty :(</EmptyTitle>
                    <EmptyDescription>
                        This repository doesn&apos;t have any commits yet. Once
                        you push your first commit, you&apos;ll be able to
                        browse the repository contents here.
                    </EmptyDescription>
                </EmptyHeader>
                <EmptyContent>
                    <Button>
                        <IconFolderPlus /> Create empty branch
                    </Button>
                </EmptyContent>
            </Empty>
        );
    }

    return (
        <RepositoryHomepage
            ownerName={ownerName}
            repoName={repoName}
            currentRef={repo.defaultBranch}
        />
    );
}
