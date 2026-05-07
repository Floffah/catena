import RepositoryBranchSelectInner from "@/components/views/RepositoryHomepage/RepositoryBranchSelect/inner";
import { serverListRepositoryRefs } from "@/lib/server/repository";

export default async function RepositoryBranchSelect({
    ownerName,
    repositoryName,
    currentRef,
}: {
    ownerName: string;
    repositoryName: string;
    currentRef: string;
}) {
    const availableRefs = await serverListRepositoryRefs(
        ownerName,
        repositoryName,
        "branch",
    );

    return (
        <RepositoryBranchSelectInner
            ownerName={ownerName}
            repositoryName={repositoryName}
            currentRef={currentRef}
            availableBranches={availableRefs?.refs.map((ref) => ref.name)}
        />
    );
}
