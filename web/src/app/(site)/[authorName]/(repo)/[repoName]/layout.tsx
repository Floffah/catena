import { notFound } from "next/navigation";
import { PropsWithChildren } from "react";

import { serverGetRepository } from "@/lib/server/repository";

export default async function Layout({
    children,
    params,
}: PropsWithChildren<{
    params: Promise<{ authorName: string; repoName: string }>;
}>) {
    const { authorName, repoName } = await params;

    const repo = await serverGetRepository(authorName, repoName);

    if (!repo) {
        return notFound();
    }

    return children;
}
