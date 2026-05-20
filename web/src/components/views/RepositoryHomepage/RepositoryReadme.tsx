import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { Card, CardContent } from "@/components/ui/card";
import { serverGetRepositoryReadme } from "@/lib/server/repository";

export async function RepositoryReadme({
    ownerName,
    repoName,
    currentRef,
    path,
}: {
    ownerName: string;
    repoName: string;
    currentRef?: string;
    path?: string;
}) {
    const readme = await serverGetRepositoryReadme(
        ownerName,
        repoName,
        currentRef,
        path,
    );

    if (!readme) {
        return null;
    }

    return (
        <Card asChild>
            <article>
                <CardContent className="prose prose-sm dark:prose-invert">
                    <Markdown remarkPlugins={[remarkGfm]}>
                        {readme?.content || "No README found."}
                    </Markdown>
                </CardContent>
            </article>
        </Card>
    );
}
