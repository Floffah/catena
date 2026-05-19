import RepositoryBranchSelectInner from "@/components/views/RepositoryHomepage/RepositoryBranchSelect/inner";
import { serverListRepositoryRefs } from "@/lib/server/repository";

export default async function RepositoryBranchSelect({
    ownerName,
    repoName,
    currentRef,
}: {
    ownerName: string;
    repoName: string;
    currentRef: string;
}) {
    const availableRefs = await serverListRepositoryRefs(
        ownerName,
        repoName,
        "branch",
    );

    return (
        <RepositoryBranchSelectInner
            ownerName={ownerName}
            repoName={repoName}
            currentRef={currentRef}
            availableBranches={availableRefs?.refs.map((ref) => ref.name)}
        />
    );
}
