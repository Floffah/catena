"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ComponentProps } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export default function SettingsNavLinkButton({
    href,
    children,
    className,
    ...props
}: Omit<ComponentProps<typeof Button>, "variant"> & { href: string }) {
    const pathname = usePathname();

    return (
        <Button
            variant={pathname === href ? "outline" : "ghost"}
            asChild
            className={cn("justify-start", className)}
            {...props}
        >
            <Link href={href}>{children}</Link>
        </Button>
    );
}
