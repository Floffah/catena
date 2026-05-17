"use client";

import { useClerk } from "@clerk/nextjs";
import { Slot } from "radix-ui";
import { ComponentProps } from "react";

export default function UserProfileDialogButton({
    asChild,
    onClick,
    startPath,
    ...props
}: ComponentProps<"button"> & { asChild?: boolean; startPath?: string }) {
    const { openUserProfile } = useClerk();

    const Comp = asChild ? Slot.Root : "button";

    return (
        <Comp
            onClick={(e) => {
                onClick?.(e);

                if (!e.defaultPrevented && !e.isPropagationStopped()) {
                    openUserProfile({ __experimental_startPath: startPath });
                }
            }}
            {...props}
        />
    );
}
