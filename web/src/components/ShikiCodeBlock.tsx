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

enum SupportedLanguages {
    JavaScript = "javascript",
    TypeScript = "typescript",
    Python = "python",
    Go = "go",
    Java = "java",
    Ruby = "ruby",
    PHP = "php",
    Rust = "rust",
    Cpp = "cpp",
    CSharp = "csharp",
    HTML = "html",
    CSS = "css",
    JSON = "json",
    Markdown = "markdown",
    YAML = "yaml",
    Dotenv = "dotenv",
    Bash = "bash",
    Docker = "docker",
    Makefile = "makefile",
    CMake = "cmake",
    Just = "just",
    INI = "ini",
}

const fileNamesToLanguages: Record<string, SupportedLanguages> = {
    Dockerfile: SupportedLanguages.Docker,
    Makefile: SupportedLanguages.Makefile,
    "CMakeLists.txt": SupportedLanguages.CMake,
    Gemfile: SupportedLanguages.Ruby,
    Rakefile: SupportedLanguages.Ruby,
    Vagrantfile: SupportedLanguages.Ruby,
    Podfile: SupportedLanguages.Ruby,
    Brewfile: SupportedLanguages.Ruby,
    Justfile: SupportedLanguages.Just,
    "go.mod": SupportedLanguages.Go,
    "go.sum": SupportedLanguages.Go,
    "bun.lock": SupportedLanguages.JSON,
    "yarn.lock": SupportedLanguages.JSON,
    "tern.conf": SupportedLanguages.INI,
    ".env.example": SupportedLanguages.Dotenv,
};

export function fileNameToLanguage(fileName: string): BundledLanguage | string {
    if (fileName in fileNamesToLanguages) {
        return fileNamesToLanguages[fileName];
    }

    const ext = fileName.split(".").pop()?.toLowerCase();

    if (!ext) {
        return SupportedLanguages.TypeScript;
    }

    // if (ext in bundledLanguages) {
    //     return ext as BundledLanguage;
    // }

    switch (ext) {
        case "js":
        case "mjs":
        case "cjs":
        case "jsx":
            return SupportedLanguages.JavaScript;
        case "ts":
        case "mts":
        case "cts":
        case "tsx":
            return SupportedLanguages.TypeScript;
        case "py":
            return SupportedLanguages.Python;
        case "go":
            return SupportedLanguages.Go;
        case "java":
            return SupportedLanguages.Java;
        case "rb":
            return SupportedLanguages.Ruby;
        case "php":
            return SupportedLanguages.PHP;
        case "rs":
            return SupportedLanguages.Rust;
        case "cpp":
        case "cc":
        case "cxx":
        case "c":
            return SupportedLanguages.Cpp;
        case "cs":
            return SupportedLanguages.CSharp;
        case "html":
        case "htm":
            return SupportedLanguages.HTML;
        case "css":
            return SupportedLanguages.CSS;
        case "json":
            return SupportedLanguages.JSON;
        case "md":
            return SupportedLanguages.Markdown;
        case "yml":
        case "yaml":
            return SupportedLanguages.YAML;
        case "env":
            return SupportedLanguages.Dotenv;
        case "sh":
            return SupportedLanguages.Bash;
        case "sql":
            return "sql";
        default:
            if (ext && ext in bundledLanguages) {
                return ext as BundledLanguage;
            }
            return SupportedLanguages.TypeScript;
    }
}

const highlighter = await createHighlighter({
    themes: [
        theme as unknown as ThemeInput,
        lightTheme as unknown as ThemeInput,
    ],
    langs: Object.values(SupportedLanguages),
});

export default async function ShikiCodeBlock({
    children,
    lang,
    className,
    ...props
}: ComponentProps<"pre"> & {
    children: string;
    lang: BundledLanguage | string;
}) {
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
