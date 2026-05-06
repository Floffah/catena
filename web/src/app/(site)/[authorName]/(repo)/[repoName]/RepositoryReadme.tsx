import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { Card, CardContent } from "@/components/ui/card";
import {
    serverGetRepository,
    serverGetRepositoryReadme,
} from "@/lib/server/repository";

export async function RepositoryReadme({
    ownerName,
    repositoryName,
}: {
    ownerName: string;
    repositoryName: string;
}) {
    const readme = await serverGetRepositoryReadme(ownerName, repositoryName);

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
