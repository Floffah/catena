import { IconFile } from "@tabler/icons-react";

import ReadmeMarkdown from "@/components/ReadmeMarkdown";
import ShikiCodeBlock, {
    fileNameToLanguage,
} from "@/components/ShikiCodeBlock";
import { serverGetRepositoryFile } from "@/lib/server/repository";

function isReadmePath(path: string) {
    const name = path.split("/").at(-1)?.toLowerCase();

    return name === "readme.md" || name === "readme";
}

export default async function RepositoryFileViewer({
    ownerName,
    repoName,
    branch,
    path,
}: {
    ownerName: string;
    repoName: string;
    branch: string;
    path: string;
}) {
    const file = await serverGetRepositoryFile(
        ownerName,
        repoName,
        branch,
        path,
    );

    if (!file) {
        return null;
    }

    const shouldRenderMarkdown = isReadmePath(file.path);

    return (
        <main className="flex flex-col gap-4">
            <header className="flex items-center gap-2 rounded-lg bg-card px-3 py-2 text-sm ring-1 ring-foreground/10">
                <IconFile className="size-4 text-muted-foreground" />
                <span className="min-w-0 flex-1 truncate font-medium">
                    {file.path}
                </span>
                <span className="text-xs text-muted-foreground">
                    {file.size.toLocaleString()} bytes
                </span>
            </header>
            {shouldRenderMarkdown && (
                <article className="rounded-lg bg-card p-4 ring-1 ring-foreground/10">
                    <ReadmeMarkdown content={file.content} />
                </article>
            )}
            {!shouldRenderMarkdown && (
                <div className="overflow-hidden rounded-lg bg-card ring-1 ring-foreground/10">
                    <ShikiCodeBlock
                        lang={fileNameToLanguage(file.name)}
                        className="overflow-x-auto p-4 text-xs leading-relaxed"
                    >
                        {file.content}
                    </ShikiCodeBlock>
                </div>
            )}
        </main>
    );
}
