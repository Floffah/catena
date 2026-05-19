"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ComponentProps } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export default function RepoNavLink({
    href,
    className,
    children,
    ...props
}: ComponentProps<typeof Button> & { href: string }) {
    const pathname = usePathname();

    return (
        <Button
            variant="ghost"
            size="lg"
            asChild
            data-active={pathname === href}
            className={cn(
                "relative data-active:after:absolute data-active:after:inset-0 data-active:after:-bottom-2.5 data-active:after:border-b data-active:after:border-primary/80",
                className,
            )}
            {...props}
        >
            <Link href={href}>{children}</Link>
        </Button>
    );
}
