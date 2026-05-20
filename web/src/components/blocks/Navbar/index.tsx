"use server";

import Link from "next/link";

import NewRepoButton from "@/components/blocks/Navbar/NewRepoButton";
import ProfileDropdownButton from "@/components/blocks/Navbar/ProfileDropdownButton";

export default async function Navbar() {
    return (
        <nav className="w-full border-b peer-data-displaces-nav:border-b-transparent next-peer-data-displaces-nav:border-b-transparent">
            <div className="container mx-auto flex items-center justify-between px-8 py-4">
                <h1 className="text-lg font-bold">
                    <Link href="/home">Catena</Link>
                </h1>
                <div className="flex items-center gap-4">
                    <NewRepoButton />

                    <ProfileDropdownButton />
                </div>
            </div>
        </nav>
    );
}
