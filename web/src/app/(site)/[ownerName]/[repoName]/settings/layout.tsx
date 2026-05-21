import { IconSettings } from "@tabler/icons-react";
import { PropsWithChildren } from "react";

import SettingsNavLinkButton from "@/components/SettingsNavLinkButton";

export default async function Layout({
    children,
    params,
}: PropsWithChildren<{
    params: Promise<{ ownerName: string; repoName: string }>;
}>) {
    const { ownerName, repoName } = await params;

    return (
        <div className="flex flex-1 gap-10">
            <nav className="flex max-w-48 shrink-0 flex-col gap-4">
                <section className="flex flex-col gap-1">
                    <h2 className="text-sm font-semibold text-muted-foreground">
                        Repository
                    </h2>

                    <SettingsNavLinkButton
                        href={`/${ownerName}/${repoName}/settings`}
                    >
                        <IconSettings className="size-4" />
                        General
                    </SettingsNavLinkButton>
                </section>
            </nav>

            {children}
        </div>
    );
}
