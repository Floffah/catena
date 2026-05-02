"use client";

import { useStore } from "@tanstack/react-form";
import { useId } from "react";

import { Field, FieldError, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { useFieldContext } from "@/hooks/useAppForm";

interface InputFieldProps {
    label: string;
}

export default function InputField({ label }: InputFieldProps) {
    const id = useId();
    const field = useFieldContext<string>();

    const errors = useStore(field.store, (state) => state.meta.errors);

    return (
        <Field>
            <FieldLabel htmlFor={id}>{label}</FieldLabel>
            <Input
                id={id}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                onBlur={field.handleBlur}
                className={errors.length > 0 ? "border-destructive" : ""}
            />
            {errors.length > 0 && (
                <FieldError>
                    {errors.map((e) => e.message ?? String(e))}
                </FieldError>
            )}
        </Field>
    );
}
