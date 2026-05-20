"use client";

import { useStore } from "@tanstack/react-form";
import { useId } from "react";

import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Textarea } from "@/components/ui/textarea";
import { useFieldContext } from "@/hooks/useAppForm";

interface TextareaFieldProps {
    label: string;
}

export default function TextareaField({ label }: TextareaFieldProps) {
    const id = useId();
    const field = useFieldContext<string>();

    const errors = useStore(field.store, (state) => state.meta.errors);

    return (
        <Field>
            <FieldLabel htmlFor={id}>{label}</FieldLabel>
            <Textarea
                id={id}
                className={errors.length > 0 ? "border-destructive" : ""}
                value={field.state.value}
                onBlur={field.handleBlur}
                onChange={(e) => field.handleChange(e.target.value)}
            />
            {errors.length > 0 && (
                <FieldError>
                    {errors.map((e) => e.message ?? String(e))}
                </FieldError>
            )}
        </Field>
    );
}
