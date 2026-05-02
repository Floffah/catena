"use client";

import { useRouter } from "next/navigation";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { FieldError } from "@/components/ui/field";
import { useAppForm, useFormContext } from "@/hooks/useAppForm";
import useUser from "@/hooks/useUser";
import { $api } from "@/lib/api";

const formSchema = z.object({
    name: z.string().min(1, "Name is required"),
    description: z.string().optional(),
});

export default function Page() {
    const router = useRouter();
    const user = useUser();

    const createRepoMutation = $api.useMutation("post", "/v1/repositories", {
        onSuccess: (data) => {
            router.push(`/${data.ownerName}/${data.name}`);
        },
    });

    const form = useAppForm({
        defaultValues: {
            name: "",
        } as z.infer<typeof formSchema>,
        validators: {
            onSubmit: formSchema,
        },
        onSubmit: async ({ value }) => {
            createRepoMutation.mutate({ body: value });
        },
    });

    return (
        <form
            onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
            }}
            className="mx-auto flex w-full max-w-lg flex-col gap-4 p-4"
        >
            <div className="flex flex-col gap-1">
                <form.AppField
                    name="name"
                    children={(field) => (
                        <field.InputField label="Repository Name" />
                    )}
                />

                <form.Subscribe
                    selector={(state) => state.values.name}
                    children={(name) => (
                        <p className="text-xs text-muted-foreground">
                            This will be: {user?.name}/
                            {name || "repository-name"}
                        </p>
                    )}
                />
            </div>

            <form.AppField
                name="description"
                children={(field) => (
                    <field.InputField label="Description (optional)" />
                )}
            />

            <div className="flex items-center gap-1">
                {createRepoMutation.error && (
                    <FieldError>
                        {createRepoMutation.error.error ||
                            "An error occurred while creating the repository."}
                    </FieldError>
                )}

                <Button
                    type="submit"
                    className="ml-auto"
                    disabled={createRepoMutation.isPending}
                >
                    Create Repository
                </Button>
            </div>
        </form>
    );
}
