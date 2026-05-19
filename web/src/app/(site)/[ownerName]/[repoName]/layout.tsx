import { PropsWithChildren } from "react";

import RepoLayout from "@/components/layouts/RepoLayout";

export default async function Layout({
    params,
    children,
}: PropsWithChildren<{
    params: Promise<{ ownerName: string; repoName: string }>;
}>) {
    const { ownerName, repoName } = await params;

    return (
        <RepoLayout ownerName={ownerName} repoName={repoName}>
            {children}
        </RepoLayout>
    );
}
