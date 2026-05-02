import { createFormHook, createFormHookContexts } from "@tanstack/react-form";

import InputField from "@/components/form/InputField";

export const { fieldContext, formContext, useFormContext, useFieldContext } =
    createFormHookContexts();

export const { useAppForm } = createFormHook({
    fieldComponents: {
        InputField,
    },
    formComponents: {},
    fieldContext,
    formContext,
});
