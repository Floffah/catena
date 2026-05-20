import { IconUserCog, IconUserEdit } from "@tabler/icons-react";
import Link from "next/link";
import { notFound } from "next/navigation";

import UserAvatar from "@/components/UserAvatar";
import UserProfileDialogButton from "@/components/UserProfileDialogButton";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { serverGetUserForName } from "@/lib/server/users";

export default async function Page({
    params,
}: {
    params: Promise<{ ownerName: string }>;
}) {
    const { ownerName } = await params;

    const user = await serverGetUserForName(ownerName);

    if (!user) {
        return notFound();
    }

    return (
        <main className="container mx-auto flex flex-1 flex-col gap-6 p-4">
            <section className="grid gap-6 md:grid-cols-[16rem_minmax(0,1fr)] md:items-start">
                <aside className="flex justify-center md:justify-start">
                    <UserAvatar user={user} className="size-48 md:size-64" />
                </aside>

                <div className="flex min-w-0 flex-col gap-6">
                    <header className="flex flex-col gap-2 text-center md:text-left">
                        <h1 className="font-heading text-4xl font-bold tracking-tight md:text-5xl">
                            {user.name}
                        </h1>
                        <p className="max-w-2xl text-sm/relaxed text-muted-foreground md:text-base/relaxed">
                            Lorem ipsum dolor sit amet, consectetur adipiscing
                            elit. Integer vitae justo sed nibh cursus
                            ullamcorper.
                        </p>
                    </header>

                    <div className="flex w-full items-center gap-2">
                        <Button asChild variant="secondary">
                            <Link
                                href={`/${user.name}/settings`}
                                className="flex flex-1 items-center gap-1"
                            >
                                <IconUserEdit /> Settings
                            </Link>
                        </Button>
                        <UserProfileDialogButton asChild>
                            <Button variant="secondary" className="flex-1">
                                <IconUserCog /> Manage Account
                            </Button>
                        </UserProfileDialogButton>
                    </div>

                    <Card asChild>
                        <article>
                            <CardHeader>
                                <CardTitle>README</CardTitle>
                            </CardHeader>
                            <CardContent className="prose prose-sm max-w-none dark:prose-invert">
                                <p>
                                    Lorem ipsum dolor sit amet, consectetur
                                    adipiscing elit. Praesent eget nibh non
                                    neque fermentum luctus. Curabitur ac
                                    facilisis turpis, sed tincidunt neque.
                                </p>
                            </CardContent>
                        </article>
                    </Card>
                </div>
            </section>

            <section className="grid gap-4">
                <Card>
                    <CardHeader>
                        <CardTitle>Projects</CardTitle>
                    </CardHeader>
                    <CardContent className="text-sm text-muted-foreground">
                        Projects will appear here.
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>Contribution history</CardTitle>
                    </CardHeader>
                    <CardContent className="text-sm text-muted-foreground">
                        Contribution history will appear here.
                    </CardContent>
                </Card>
            </section>
        </main>
    );
}
