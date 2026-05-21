"use client";

import { useRouter } from "next/navigation";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { useAppForm } from "@/hooks/useAppForm";
import { $api } from "@/lib/api";
import { SchemaRepository, SchemaRepositoryVisibility } from "@/types/api";

const formSchema = z.object({
    description: z.string(),
    defaultBranch: z.string().min(1, "Default branch is required"),
    visibility: z.enum(["private", "public"]),
});

export default function RepositorySettingsForm({
    repository,
    branchNames,
}: {
    repository: SchemaRepository;
    branchNames: string[];
}) {
    const router = useRouter();

    const updateRepositoryMutation = $api.useMutation(
        "patch",
        "/v1/repositories/{owner}/{repository}",
        {
            onSuccess: () => {
                router.refresh();
            },
        },
    );

    const form = useAppForm({
        defaultValues: {
            description: repository.description ?? "",
            defaultBranch: repository.defaultBranch,
            visibility: repository.visibility,
        },
        validators: {
            onSubmit: formSchema,
        },
        onSubmit: async ({ value }) => {
            await updateRepositoryMutation.mutateAsync({
                params: {
                    path: {
                        owner: repository.ownerName,
                        repository: repository.name,
                    },
                },
                body: {
                    description: value.description,
                    defaultBranch: value.defaultBranch,
                    visibility: value.visibility as SchemaRepositoryVisibility,
                },
            });
        },
    });

    return (
        <main className="flex flex-1">
            <form
                className="flex w-full max-w-lg flex-1 flex-col gap-4"
                onSubmit={(e) => {
                    e.preventDefault();
                    form.handleSubmit();
                }}
            >
                <div>
                    <h2 className="text-xl font-bold">Repository Settings</h2>
                    <p className="text-sm text-muted-foreground">
                        Update the basic settings for {repository.ownerName}/
                        {repository.name}.
                    </p>
                </div>

                <form.AppField name="description">
                    {(field) => (
                        <field.TextareaField label="Description (optional)" />
                    )}
                </form.AppField>

                <form.Field name="defaultBranch">
                    {(field) => (
                        <Field>
                            <FieldLabel htmlFor={field.name}>
                                Default Branch
                            </FieldLabel>
                            <select
                                id={field.name}
                                className="h-8 w-full rounded-md border border-input bg-input/20 px-2 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 disabled:cursor-not-allowed disabled:opacity-50 md:text-xs/relaxed dark:bg-input/30"
                                value={field.state.value}
                                onBlur={field.handleBlur}
                                onChange={(e) =>
                                    field.handleChange(e.target.value)
                                }
                            >
                                {branchNames.map((branchName) => (
                                    <option key={branchName} value={branchName}>
                                        {branchName}
                                    </option>
                                ))}
                            </select>
                        </Field>
                    )}
                </form.Field>

                <form.Field name="visibility">
                    {(field) => (
                        <Field>
                            <FieldLabel htmlFor={field.name}>
                                Visibility
                            </FieldLabel>
                            <select
                                id={field.name}
                                className="h-8 w-full rounded-md border border-input bg-input/20 px-2 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 disabled:cursor-not-allowed disabled:opacity-50 md:text-xs/relaxed dark:bg-input/30"
                                value={field.state.value}
                                onBlur={field.handleBlur}
                                onChange={(e) =>
                                    field.handleChange(
                                        e.target
                                            .value as SchemaRepositoryVisibility,
                                    )
                                }
                            >
                                <option value="private">Private</option>
                                <option value="public">Public</option>
                            </select>
                        </Field>
                    )}
                </form.Field>

                <div className="flex items-center gap-2">
                    {updateRepositoryMutation.error && (
                        <FieldError>
                            {updateRepositoryMutation.error.error ||
                                "An error occurred while updating the repository."}
                        </FieldError>
                    )}
                    {updateRepositoryMutation.isSuccess && (
                        <p className="text-sm text-success">
                            Repository updated successfully!
                        </p>
                    )}

                    <Button
                        className="ml-auto"
                        disabled={updateRepositoryMutation.isPending}
                        type="submit"
                    >
                        Save Changes
                    </Button>
                </div>
            </form>
        </main>
    );
}
