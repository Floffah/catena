"use client";

import { IconPlus } from "@tabler/icons-react";
import Link from "next/link";

import { Button } from "@/components/ui/button";

export default function Navbar() {
    return (
        <nav className="w-full border-b">
            <div className="container mx-auto flex items-center justify-between px-8 py-4">
                <h1 className="text-lg font-bold">Catena</h1>
                <div>
                    <Button asChild variant="secondary">
                        <Link href="/new/repository">
                            <IconPlus />
                            New Repository
                        </Link>
                    </Button>
                </div>
            </div>
        </nav>
    );
}
