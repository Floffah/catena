export const routerPushCalls: string[] = [];
export let routerRefreshCalls = 0;
export let mockParams: Record<string, string | string[]> = {};

export function resetMockNextNavigation() {
    routerPushCalls.length = 0;
    routerRefreshCalls = 0;
    mockParams = {};
}

export function setMockParams(params: Record<string, string | string[]>) {
    mockParams = params;
}

export const mockNextNavigation = {
    useParams() {
        return mockParams;
    },
    useRouter() {
        return {
            push: (href: string) => {
                routerPushCalls.push(href);
            },
            refresh: () => {
                routerRefreshCalls += 1;
            },
        };
    },
};
