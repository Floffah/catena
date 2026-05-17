"use client";

import { IconExternalLink } from "@tabler/icons-react";
import { z } from "zod";

import UserProfileDialogButton from "@/components/UserProfileDialogButton";
import { Button } from "@/components/ui/button";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { useAppForm } from "@/hooks/useAppForm";
import useUser from "@/hooks/useUser";
import { $api } from "@/lib/api";

const formSchema = z.object({
    displayName: z.string().min(1, "Display name is required"),
});

export default function ProfileSettings() {
    const user = useUser();

    const updateUserMutation = $api.useMutation("patch", "/v1/user");

    const form = useAppForm({
        defaultValues: {
            displayName: user?.displayName || "",
        },
        validators: {
            onSubmit: formSchema,
        },
        onSubmit: async ({ value }) => {
            await updateUserMutation.mutateAsync({
                body: value,
            });
        },
    });

    return (
        <main className="flex flex-1">
            <form
                className="flex w-full max-w-64 flex-1 flex-col gap-4"
                onSubmit={(e) => {
                    e.preventDefault();
                    form.handleSubmit();
                }}
            >
                <h2 className="text-xl font-bold">My Profile</h2>

                <form.AppField
                    name="displayName"
                    children={(field) => (
                        <field.InputField label="Display Name" />
                    )}
                />

                <Field>
                    <FieldLabel>Username</FieldLabel>
                    <div className="flex items-center gap-2">
                        <Input disabled value={user.name} />
                        <Button variant="outline" asChild>
                            <UserProfileDialogButton>
                                <IconExternalLink className="size-4" />
                                Edit
                            </UserProfileDialogButton>
                        </Button>
                    </div>
                </Field>

                <Field>
                    <FieldLabel>Email</FieldLabel>
                    <div className="flex items-center gap-2">
                        <Input disabled value={user.email || "Not provided"} />
                        <Button variant="outline" asChild>
                            <UserProfileDialogButton>
                                <IconExternalLink className="size-4" />
                                Edit
                            </UserProfileDialogButton>
                        </Button>
                    </div>
                </Field>

                <Button type="submit" disabled={updateUserMutation.isPending}>
                    Save Changes
                </Button>
            </form>
        </main>
    );
}
