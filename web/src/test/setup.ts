import { afterAll, afterEach, beforeAll, mock } from "bun:test";
import { Window } from "happy-dom";

import { mockNextNavigation, resetMockNextNavigation } from "@/test/navigation";
import { server } from "@/test/server";

process.env.NEXT_PUBLIC_CATENA_INSTANCE_URL = "http://catena.test";

mock.module("next/navigation", () => mockNextNavigation);

const window = new Window({
    url: "http://localhost",
});

window.SyntaxError = SyntaxError;

const globals = globalThis as Record<string, unknown>;

globals.window = window;
globals.document = window.document;
globals.navigator = window.navigator;
globals.HTMLElement = window.HTMLElement;
globals.HTMLButtonElement = window.HTMLButtonElement;
globals.HTMLInputElement = window.HTMLInputElement;
globals.HTMLTextAreaElement = window.HTMLTextAreaElement;
globals.Event = window.Event;
globals.MouseEvent = window.MouseEvent;
globals.KeyboardEvent = window.KeyboardEvent;
globals.CustomEvent = window.CustomEvent;
globals.Node = window.Node;
globals.NodeFilter = window.NodeFilter;
globals.Text = window.Text;
globals.Element = window.Element;
globals.MutationObserver = window.MutationObserver;
globals.ResizeObserver = window.ResizeObserver;
globals.DOMParser = window.DOMParser;
globals.getComputedStyle = window.getComputedStyle.bind(window);
globals.requestAnimationFrame = window.requestAnimationFrame.bind(window);
globals.cancelAnimationFrame = window.cancelAnimationFrame.bind(window);
globals.SyntaxError = SyntaxError;
globals.IS_REACT_ACT_ENVIRONMENT = true;

const { cleanup } = await import("@testing-library/react");

beforeAll(() => {
    server.listen({
        onUnhandledRequest: "error",
    });
});

afterEach(() => {
    cleanup();
    resetMockNextNavigation();
    server.resetHandlers();
});

afterAll(() => {
    server.close();
});
