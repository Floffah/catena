import { ClerkProvider } from "@clerk/nextjs";
import { shadcn } from "@clerk/ui/themes";
import type { Metadata } from "next";
import { ThemeProvider } from "next-themes";
import { Geist, JetBrains_Mono } from "next/font/google";
import { PropsWithChildren } from "react";

import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import AuthProvider from "@/providers/AuthProvider";
import QueryClientProvider from "@/providers/QueryClientProvider";

import "./globals.css";

const sansFont = Geist({ subsets: ["latin"], variable: "--font-sans" });

const monoFont = JetBrains_Mono({
    variable: "--font-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "Catena",
    description: "Where the next generation of open source software is built.",
    openGraph: {
        title: "Catena",
        description:
            "Where the next generation of open source software is built.",
        url: "https://oncatena.com",
        siteName: "Catena",
    },
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
                    <TooltipProvider>
                        <ClerkProvider appearance={{ theme: shadcn }}>
                            <QueryClientProvider>
                                <AuthProvider>
                                    {children}

                                    <Toaster />
                                </AuthProvider>
                            </QueryClientProvider>
                        </ClerkProvider>
                    </TooltipProvider>
                </ThemeProvider>
            </body>
        </html>
    );
}
