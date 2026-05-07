import { PropsWithChildren } from "react";

import RepoLayout from "@/app/(site)/[authorName]/[repoName]/RepoLayout";

export default async function Layout({
    params,
    children,
}: PropsWithChildren<{
    params: Promise<{ authorName: string; repoName: string }>;
}>) {
    const { authorName, repoName } = await params;

    return (
        <RepoLayout authorName={authorName} repoName={repoName}>
            {children}
        </RepoLayout>
    );
}
