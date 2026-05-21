import { Metadata, ResolvingMetadata } from "next";
import { PropsWithChildren } from "react";

import RepoLayout from "@/components/layouts/RepoLayout";
import { apiFetch } from "@/lib/api";

export async function generateMetadata(
    {
        params,
    }: {
        params: Promise<{ ownerName: string; repoName: string }>;
    },
    parent: ResolvingMetadata,
) {
    const { ownerName, repoName } = await params;
    try {
        const response = await apiFetch.GET(
            "/v1/repositories/{owner}/{repository}",
            {
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
                    },
                },
            },
        );

        if (response.data) {
            return {
                title: `${response.data.ownerName}/${response.data.name} - Catena`,
                description:
                    response.data.description ||
                    `${response.data.ownerName} builds ${response.data.name} on Catena.`,
                openGraph: {
                    title: `${response.data.ownerName}/${response.data.name} - Catena`,
                    description:
                        response.data.description ||
                        `${response.data.ownerName} builds ${response.data.name} on Catena.`,
                    url: `https://oncatena.com/${response.data.ownerName}/${response.data.name}`,
                    siteName: "Catena",
                },
            } satisfies Metadata;
        }
    } catch {}

    return await parent;
}

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
