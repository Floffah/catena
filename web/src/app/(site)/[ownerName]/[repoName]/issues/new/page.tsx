import CreateIssueForm from "@/components/views/CreateIssueForm";

export default async function Page({
    params,
}: {
    params: Promise<{ ownerName: string; repoName: string }>;
}) {
    const { ownerName, repoName } = await params;

    return (
        <main className="flex flex-col gap-4">
            <header>
                <h1 className="text-2xl font-bold">New issue</h1>
                <p className="text-sm text-muted-foreground">
                    Capture a bug, task, or idea for this repository.
                </p>
            </header>

            <CreateIssueForm ownerName={ownerName} repoName={repoName} />
        </main>
    );
}
