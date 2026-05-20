import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

import InputField from "@/components/form/InputField";
import TextareaField from "@/components/form/TextareaField";

export const { fieldContext, formContext, useFormContext, useFieldContext } =
    createFormHookContexts();

export const { useAppForm } = createFormHook({
    fieldComponents: {
        InputField,
        TextareaField,
    },
    formComponents: {},
    fieldContext,
    formContext,
});
