import { ClerkProvider } from "@clerk/nextjs";
import { shadcn } from "@clerk/ui/themes";
import type { Metadata } from "next";
import { ThemeProvider } from "next-themes";
import { Geist, Geist_Mono } from "next/font/google";
import { PropsWithChildren } from "react";

import { SidebarProvider } from "@/components/ui/sidebar";
import { cn } from "@/lib/utils";
import AuthProvider from "@/providers/AuthProvider";
import QueryClientProvider from "@/providers/QueryClientProvider";

import "./globals.css";

const sansFont = Geist({ subsets: ["latin"], variable: "--font-sans" });

const monoFont = Geist_Mono({
    variable: "--font-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "Catena",
    description: "Next generation git server",
};

export default function RootLayout({ children }: PropsWithChildren) {
    return (
        <html
            lang="en"
            className={cn(
                "h-full",
                "antialiased",
                sansFont.variable,
                monoFont.variable,
                "font-sans",
            )}
            suppressHydrationWarning
        >
            <body className="flex min-h-full flex-col" suppressHydrationWarning>
                <ThemeProvider
                    attribute="class"
                    defaultTheme="system"
                    enableColorScheme
                    disableTransitionOnChange
                    enableSystem
                >
                    <ClerkProvider appearance={{ theme: shadcn }}>
                        <QueryClientProvider>
                            <AuthProvider>
                                {children}

                                <SidebarProvider />
                            </AuthProvider>
                        </QueryClientProvider>
                    </ClerkProvider>
                </ThemeProvider>
            </body>
        </html>
    );
}
