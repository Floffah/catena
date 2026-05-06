"use client";

import { useParams } from "next/navigation";
import { PropsWithChildren, createContext, useContext } from "react";

import { $api } from "@/lib/api";

import { SchemaRepository } from "../../types/api";

interface RepositoryContextValue extends SchemaRepository {}

const RepositoryContext = createContext<RepositoryContextValue>(null!);

export const useRepo = () => {
    const context = useContext(RepositoryContext);

    if (!context) {
        throw new Error("useRepo must be used within a RepositoryProvider");
    }

    return context;
};

export default function RepositoryProvider({ children }: PropsWithChildren) {
    const params = useParams();

    const repositoryQuery = $api.useQuery(
        "get",
        "/v1/repositories/{owner}/{repository}",
        {
            params: {
                path: {
                    owner: params.authorName as string,
                    repository: params.repoName as string,
                },
            },
        },
        {
            meta: {
                refetchOnAuth: true,
            },
        },
    );

    return (
        <RepositoryContext.Provider value={repositoryQuery.data!}>
            {children}
        </RepositoryContext.Provider>
    );
}
