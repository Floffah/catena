import { IconUser } from "@tabler/icons-react";
import { ComponentProps } from "react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Spinner } from "@/components/ui/spinner";
import { SchemaUser } from "@/types/api";

export default function UserAvatar({
    user: { name, avatarUrl },
    ...props
}: ComponentProps<typeof Avatar> & {
    user: Partial<SchemaUser> & { name: string };
}) {
    return (
        <Avatar {...props}>
            <AvatarFallback>
                {name[0].toUpperCase() ?? <IconUser className="size-4" />}
            </AvatarFallback>
            {avatarUrl && <AvatarImage src={avatarUrl} alt={name} />}
        </Avatar>
    );
}

export function UserAvatarFallback({
    loading,
    ...props
}: ComponentProps<typeof Avatar> & { loading?: boolean }) {
    let icon = <IconUser className="size-6/12" />;

    if (loading) {
        icon = <Spinner className="size-6/12" />;
    }

    return (
        <Avatar {...props}>
            <AvatarFallback>{icon}</AvatarFallback>
        </Avatar>
    );
}
