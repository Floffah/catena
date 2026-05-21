import { toJsxRuntime } from "hast-util-to-jsx-runtime";
import { ComponentProps, Fragment } from "react";
import { JSX, jsx, jsxs } from "react/jsx-runtime";
import { ThemeInput, createHighlighter } from "shiki";
import { BundledLanguage } from "shiki/bundle/web";

import { cn } from "@/lib/utils";
import lightTheme from "@/public/code-theme-light.json";
import theme from "@/public/code-theme.json";

export default async function ShikiCodeBlock({
    children,
    lang,
    className,
    ...props
}: ComponentProps<"pre"> & {
    children: string;
    lang: BundledLanguage | string;
}) {
    const highlighter = await createHighlighter({
        themes: [
            theme as unknown as ThemeInput,
            lightTheme as unknown as ThemeInput,
        ],
        langs: ["go"],
    });
    const out = highlighter.codeToHast(children, {
        lang: lang,
        themes: {
            light: "catena-light",
            dark: "catena-dark",
        },
    });

    return toJsxRuntime(out, {
        Fragment,
        jsx,
        jsxs,
        components: {
            // your custom `pre` element
            pre: (p) => (
                <pre
                    data-codeblock
                    {...p}
                    {...props}
                    className={cn("line-numbers", p.className, className)}
                />
            ),
        },
    }) as JSX.Element;
}
