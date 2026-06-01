import { IconUserCog, IconUserEdit } from "@tabler/icons-react";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Suspense } from "react";

import UserAvatar from "@/components/UserAvatar";
import UserProfileDialogButton from "@/components/UserProfileDialogButton";
import { Button } from "@/components/ui/button";
import ProfileReadme, {
    ProfileReadmeSkeleton,
} from "@/components/views/ProfileReadme";
import ProfileRepositories, {
    ProfileRepositoriesSkeleton,
} from "@/components/views/ProfileRepositories";
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
                            {user.displayName ?? user.name}
                        </h1>
                        {user.name && user.displayName && (
                            <p className="-mt-2 text-sm text-muted-foreground">
                                @{user.name}
                            </p>
                        )}
                        {user.description && (
                            <p className="max-w-2xl text-sm/relaxed whitespace-pre-line text-muted-foreground md:text-base/relaxed">
                                {user.description}
                            </p>
                        )}
                    </header>

                    <div className="flex w-full items-center gap-2">
                        <Button asChild variant="secondary">
                            <Link
                                href={`/settings`}
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

                    <Suspense fallback={<ProfileReadmeSkeleton />}>
                        <ProfileReadme ownerName={ownerName} />
                    </Suspense>
                </div>
            </section>

            <Suspense fallback={<ProfileRepositoriesSkeleton />}>
                <ProfileRepositories ownerName={ownerName} />
            </Suspense>
        </main>
    );
}
