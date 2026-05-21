import { IconCode, IconGitBranch, IconVersions } from "@tabler/icons-react";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Suspense } from "react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import RepositoryBranchSelect from "@/components/views/RepositoryHomepage/RepositoryBranchSelect";
import { serverGetRepository } from "@/lib/server/repository";

export async function RepositoryHeader({
    ownerName,
    repoName,
    currentRef,
}: {
    ownerName: string;
    repoName: string;
    currentRef: string;
}) {
    const repo = await serverGetRepository(ownerName, repoName);

    if (!repo) {
        return notFound();
    }

    return (
        <div className="flex flex-col gap-1">
            <div className="flex items-center gap-4">
                <div className="flex items-center gap-1">
                    <Suspense fallback={currentRef}>
                        <RepositoryBranchSelect
                            ownerName={ownerName}
                            repoName={repoName}
                            currentRef={currentRef}
                        />
                    </Suspense>
                    <Separator orientation="vertical" className="mx-2" />
                    <Button variant="outline" asChild>
                        <Link href={`/${ownerName}/${repoName}/refs/branches`}>
                            <IconGitBranch />
                            Branches
                        </Link>
                    </Button>
                    <Button variant="outline" asChild>
                        <Link href={`/${ownerName}/${repoName}/refs/branches`}>
                            <IconVersions />
                            Tags
                        </Link>
                    </Button>
                    <Button variant="outline" asChild>
                        <a href="#browse">
                            <IconCode />
                            Browse
                        </a>
                    </Button>
                </div>
            </div>

            {/*<p className="flex items-center gap-1 text-sm text-muted-foreground underline">*/}
            {/*    <IconCornerDownRight className="size-4" />*/}
            {/*    <a href="#browse">*/}
            {/*        Scroll down to browse the repository contents*/}
            {/*    </a>*/}
            {/*</p>*/}
        </div>
    );
}
