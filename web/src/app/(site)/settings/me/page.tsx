import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { notFound } from "next/navigation";

import ProfileSettings from "@/components/views/ProfileSettings";
import { $api } from "@/lib/api";
import { queryClient } from "@/lib/queryClient";
import { serverGetCurrentUser } from "@/lib/server/users";

export default async function Page() {
    const user = await serverGetCurrentUser();

    if (!user) {
        return notFound();
    }

    const { queryKey } = $api.queryOptions("get", "/v1/user");
    queryClient.setQueryData(queryKey, user);

    return (
        <HydrationBoundary state={dehydrate(queryClient)}>
            <ProfileSettings />
        </HydrationBoundary>
    );
}
