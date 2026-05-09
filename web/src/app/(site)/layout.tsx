import { PropsWithChildren } from "react";

import Navbar from "../../components/blocks/Navbar";

export default function Layout({ children }: PropsWithChildren) {
    return (
        <div className="flex flex-1 flex-col gap-2">
            <Navbar />

            {children}
        </div>
    );
}
