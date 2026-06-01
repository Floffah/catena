import ReadmeMarkdown from "@/components/ReadmeMarkdown";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
    serverGetRepository,
    serverGetRepositoryReadme,
} from "@/lib/server/repository";

export function ProfileReadmeSkeleton() {
    return (
        <Card asChild>
            <article aria-label="Loading README">
                <CardHeader>
                    <Skeleton className="h-4 w-20" />
                </CardHeader>
                <CardContent className="space-y-3">
                    <Skeleton className="h-4 w-3/4" />
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-11/12" />
                    <Skeleton className="h-4 w-2/3" />
                </CardContent>
            </article>
        </Card>
    );
}

export default async function ProfileReadme({
    ownerName,
}: {
    ownerName: string;
}) {
    const repository = await serverGetRepository(ownerName, ownerName);

    if (!repository) {
        return null;
    }

    const readme = await serverGetRepositoryReadme(
        ownerName,
        ownerName,
        repository.defaultBranch,
    );

    if (!readme) {
        return null;
    }

    return (
        <Card asChild>
            <article>
                <CardHeader>
                    <CardTitle>README</CardTitle>
                </CardHeader>
                <CardContent>
                    <ReadmeMarkdown content={readme.content} />
                </CardContent>
            </article>
        </Card>
    );
}
