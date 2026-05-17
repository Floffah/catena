"use client";

import { IconCheck, IconCopy } from "@tabler/icons-react";
import Link from "next/link";
import { useState } from "react";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
    InputGroup,
    InputGroupAddon,
    InputGroupButton,
    InputGroupInput,
} from "@/components/ui/input-group";
import { useAppForm } from "@/hooks/useAppForm";
import { $api } from "@/lib/api";

const formSchema = z.object({
    name: z.string().min(1, "Name is required"),
    expiresAt: z.string(),
});

export default function Page() {
    const [createdToken, setCreatedToken] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);

    const createTokenMutation = $api.useMutation(
        "post",
        "/v1/git-access-tokens",
        {
            onSuccess: ({ token }) => {
                setCreatedToken(token);
            },
        },
    );

    const form = useAppForm({
        defaultValues: {
            name: "",
            expiresAt: "",
        },
        validators: {
            onSubmit: formSchema,
        },
        onSubmit: async ({ value }) => {
            createTokenMutation.mutate({
                body: {
                    name: value.name,
                    ...(value.expiresAt
                        ? { expiresAt: new Date(value.expiresAt).toISOString() }
                        : {}),
                },
            });
        },
    });

    async function copyToken() {
        if (!createdToken) {
            return;
        }

        await navigator.clipboard.writeText(createdToken);
        setCopied(true);
    }

    return (
        <>
            <main className="flex flex-1">
                <form
                    className="flex w-full max-w-sm flex-1 flex-col gap-4"
                    onSubmit={(e) => {
                        e.preventDefault();
                        form.handleSubmit();
                    }}
                >
                    <h2 className="text-xl font-bold">
                        Create Personal Access Token
                    </h2>

                    <form.AppField name="name">
                        {(field) => <field.InputField label="Token Name" />}
                    </form.AppField>

                    <form.Field name="expiresAt">
                        {(field) => (
                            <Field>
                                <FieldLabel htmlFor={field.name}>
                                    Expiration (optional)
                                </FieldLabel>
                                <Input
                                    id={field.name}
                                    type="datetime-local"
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={(e) =>
                                        field.handleChange(e.target.value)
                                    }
                                />
                            </Field>
                        )}
                    </form.Field>

                    <div className="flex items-center gap-1">
                        {createTokenMutation.error && (
                            <FieldError>
                                {createTokenMutation.error.error ||
                                    "An error occurred while creating the token."}
                            </FieldError>
                        )}

                        <Button
                            className="ml-auto"
                            disabled={createTokenMutation.isPending}
                            type="submit"
                        >
                            Create Token
                        </Button>
                    </div>
                </form>
            </main>

            <Dialog open={createdToken !== null}>
                <DialogContent showCloseButton={false}>
                    <DialogHeader>
                        <DialogTitle>Copy your token now</DialogTitle>
                        <DialogDescription>
                            This is the only time we will show this token. Copy
                            it now before continuing.
                        </DialogDescription>
                    </DialogHeader>

                    <InputGroup>
                        <InputGroupInput readOnly value={createdToken ?? ""} />
                        <InputGroupAddon align="inline-end">
                            <InputGroupButton
                                aria-label="Copy token"
                                onClick={copyToken}
                                size="icon-xs"
                            >
                                {copied ? <IconCheck /> : <IconCopy />}
                            </InputGroupButton>
                        </InputGroupAddon>
                    </InputGroup>

                    <DialogFooter>
                        <Button asChild>
                            <Link href="/settings/tokens">Okay</Link>
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}
