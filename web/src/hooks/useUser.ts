import { useAuthedQuery } from "@/lib/api";
import { useCatenaAuth } from "@/providers/AuthProvider";

export default function useUser() {
    const auth = useCatenaAuth();

    const userQuery = useAuthedQuery("get", "/v1/user");

    return {
        isLoading: auth.isLoading || userQuery.isPending,
        isAuthenticated: auth.isAuthenticated && !!userQuery.data,
        ...userQuery.data,
    };
}
