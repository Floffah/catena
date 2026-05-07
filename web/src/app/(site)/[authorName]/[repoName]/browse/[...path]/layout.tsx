import { PropsWithChildren } from "react";

import RepoLayout from "@/app/(site)/[authorName]/[repoName]/RepoLayout";

export default async function Layout({
    params,
    children,
}: PropsWithChildren<{
    params: Promise<{ authorName: string; repoName: string; path: string[] }>;
}>) {
    const { authorName, repoName, path } = await params;

    return (
        <RepoLayout authorName={authorName} repoName={repoName} path={path}>
            {children}
        </RepoLayout>
    );
}
