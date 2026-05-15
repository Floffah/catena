import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { Card, CardContent } from "@/components/ui/card";
import { serverGetRepositoryReadme } from "@/lib/server/repository";

export async function RepositoryReadme({
    ownerName,
    repositoryName,
    branch,
    path,
}: {
    ownerName: string;
    repositoryName: string;
    branch?: string;
    path?: string;
}) {
    const readme = await serverGetRepositoryReadme(
        ownerName,
        repositoryName,
        branch,
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
