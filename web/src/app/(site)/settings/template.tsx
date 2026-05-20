import { IconKey, IconSettings, IconUserEdit } from "@tabler/icons-react";
import { PropsWithChildren } from "react";

import SettingsNavLinkButton from "@/components/SettingsNavLinkButton";
import UserProfileDialogButton from "@/components/UserProfileDialogButton";
import { Button } from "@/components/ui/button";

export default function Template({ children }: PropsWithChildren) {
    return (
        <div className="container mx-auto flex flex-1 gap-10 p-4">
            <nav className="flex max-w-48 shrink-0 flex-col gap-4">
                <section className="flex flex-col gap-1">
                    <h2 className="text-sm font-semibold text-muted-foreground">
                        Account
                    </h2>

                    <SettingsNavLinkButton href="/settings/me">
                        <IconUserEdit className="size-4" />
                        My Profile
                    </SettingsNavLinkButton>
                    <Button variant="ghost" className="justify-start" asChild>
                        <UserProfileDialogButton>
                            <IconSettings className="size-4" />
                            Manage Account
                        </UserProfileDialogButton>
                    </Button>
                </section>

                <section className="flex flex-col gap-1">
                    <h2 className="text-sm font-semibold text-muted-foreground">
                        Security
                    </h2>

                    <SettingsNavLinkButton href="/settings/tokens">
                        <IconKey className="size-4" />
                        Personal Access Tokens
                    </SettingsNavLinkButton>
                </section>
            </nav>

            {children}
        </div>
    );
}
