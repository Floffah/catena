"use client";

import { useRouter } from "next/navigation";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { FieldError } from "@/components/ui/field";
import { useAppForm } from "@/hooks/useAppForm";
import { $api } from "@/lib/api";

const formSchema = z.object({
    title: z.string().refine((value) => value.trim().length > 0, {
        message: "Title is required",
    }),
    body: z.string(),
});

export default function CreateIssueForm({
    ownerName,
    repoName,
}: {
    ownerName: string;
    repoName: string;
}) {
    const router = useRouter();

    const createIssueMutation = $api.useMutation(
        "post",
        "/v1/repositories/{owner}/{repository}/issues",
        {
            onSuccess: (data) => {
                router.push(`/${ownerName}/${repoName}/issues/${data.number}`);
                router.refresh();
            },
        },
    );

    const form = useAppForm({
        defaultValues: {
            title: "",
            body: "",
        },
        validators: {
            onSubmit: formSchema,
        },
        onSubmit: async ({ value }) => {
            createIssueMutation.mutate({
                params: {
                    path: {
                        owner: ownerName,
                        repository: repoName,
                    },
                },
                body: {
                    title: value.title.trim(),
                    body: value.body.trim().length > 0 ? value.body : null,
                },
            });
        },
    });

    return (
        <form
            className="flex w-full max-w-2xl flex-col gap-4"
            onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
            }}
        >
            <form.AppField name="title">
                {(field) => <field.InputField label="Title" />}
            </form.AppField>

            <form.AppField name="body">
                {(field) => <field.TextareaField label="Body" />}
            </form.AppField>

            {createIssueMutation.error && (
                <FieldError>
                    {createIssueMutation.error.error ||
                        "An error occurred while creating the issue."}
                </FieldError>
            )}

            <div className="justify-end">
                <Button disabled={createIssueMutation.isPending}>
                    Create issue
                </Button>
            </div>
        </form>
    );
}
