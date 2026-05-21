import { IconBrandGit } from "@tabler/icons-react";
import { notFound } from "next/navigation";
import { PropsWithChildren } from "react";

import RepoNavLink from "@/components/layouts/RepoLayout/RepoNavLink";
import { serverGetRepository } from "@/lib/server/repository";
import { serverGetCurrentUser } from "@/lib/server/users";

export default async function RepoLayout({
    children,
    ownerName,
    repoName,
}: PropsWithChildren<{
    ownerName: string;
    repoName: string;
}>) {
    const [repo, user] = await Promise.all([
        serverGetRepository(ownerName, repoName),
        serverGetCurrentUser(),
    ]);

    if (!repo) {
        return notFound();
    }

    const canManageRepository = user?.name === repo.ownerName;

    return (
        <div
            className="peer -mt-4 flex flex-1 flex-col gap-4"
            data-displaces-nav="true"
        >
            <div className="flex border-b">
                <div className="container mx-auto flex items-end justify-between gap-8 px-8 pb-2">
                    <div className="flex items-center gap-1">
                        <h1 className="flex items-center gap-1 text-xl">
                            <IconBrandGit />
                            {repo.ownerName}
                            <span className="text-2xl">/</span>
                            {repo.name}
                        </h1>
                    </div>

                    <nav className="flex items-center gap-4">
                        <RepoNavLink
                            href={`/${repo.ownerName}/${repo.name}`}
                            exact
                        >
                            Repository
                        </RepoNavLink>
                        <RepoNavLink
                            href={`/${repo.ownerName}/${repo.name}/issues`}
                        >
                            Issues
                        </RepoNavLink>
                        {canManageRepository && (
                            <RepoNavLink
                                href={`/${repo.ownerName}/${repo.name}/settings`}
                            >
                                Settings
                            </RepoNavLink>
                        )}
                    </nav>
                </div>
            </div>

            <div className="container mx-auto flex flex-1 flex-col gap-4 p-4">
                {children}
            </div>
        </div>
    );
}
