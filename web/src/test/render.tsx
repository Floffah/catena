import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RenderOptions, render } from "@testing-library/react";
import { PropsWithChildren, ReactElement, act } from "react";

export function createTestQueryClient() {
    return new QueryClient({
        defaultOptions: {
            mutations: {
                retry: false,
            },
            queries: {
                retry: false,
            },
        },
    });
}

export function renderWithQueryClient(
    element: ReactElement,
    options?: RenderOptions,
) {
    const queryClient = createTestQueryClient();

    function Wrapper({ children }: PropsWithChildren) {
        return (
            <QueryClientProvider client={queryClient}>
                {children}
            </QueryClientProvider>
        );
    }

    let renderResult: ReturnType<typeof render>;

    act(() => {
        renderResult = render(element, {
            wrapper: Wrapper,
            ...options,
        });
    });

    return {
        queryClient,
        ...renderResult!,
    };
}
