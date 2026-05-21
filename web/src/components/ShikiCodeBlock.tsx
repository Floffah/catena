import { toJsxRuntime } from "hast-util-to-jsx-runtime";
import { ComponentProps, Fragment } from "react";
import { JSX, jsx, jsxs } from "react/jsx-runtime";
import {
    BundledLanguage,
    ThemeInput,
    bundledLanguages,
    createHighlighter,
} from "shiki";

import { cn } from "@/lib/utils";
import lightTheme from "@/public/code-theme-light.json";
import theme from "@/public/code-theme.json";

const fileNamesToLanguages: Record<string, BundledLanguage> = {
    Dockerfile: "docker",
    Makefile: "makefile",
    "CMakeLists.txt": "cmake",
    Gemfile: "ruby",
    Rakefile: "ruby",
    Vagrantfile: "ruby",
    Podfile: "ruby",
    Brewfile: "ruby",
    Justfile: "just",
    "go.mod": "go",
    "go.sum": "go",
    "bun.lock": "json",
    "yarn.lock": "json",
    "tern.conf": "ini",
    ".env.example": "dotenv",
};

export function fileNameToLanguage(fileName: string): BundledLanguage | string {
    if (fileName in fileNamesToLanguages) {
        return fileNamesToLanguages[fileName];
    }

    const ext = fileName.split(".").pop()?.toLowerCase();

    if (!ext) {
        return "ts";
    }

    if (ext in bundledLanguages) {
        return ext as BundledLanguage;
    }

    switch (ext) {
        case "js":
        case "mjs":
        case "cjs":
        case "jsx":
            return "javascript";
        case "ts":
        case "mts":
        case "cts":
        case "tsx":
            return "typescript";
        case "py":
            return "python";
        case "go":
            return "go";
        case "java":
            return "java";
        case "rb":
            return "ruby";
        case "php":
            return "php";
        case "rs":
            return "rust";
        case "cpp":
        case "cc":
        case "cxx":
        case "c":
            return "cpp";
        case "cs":
            return "csharp";
        case "html":
        case "htm":
            return "html";
        case "css":
            return "css";
        case "json":
            return "json";
        case "md":
            return "markdown";
        case "yml":
        case "yaml":
            return "yaml";
        case "env":
            return "dotenv";
        case "sh":
            return "bash";
        default:
            if (ext && ext in bundledLanguages) {
                return ext as BundledLanguage;
            }
            return "ts";
    }
}

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
