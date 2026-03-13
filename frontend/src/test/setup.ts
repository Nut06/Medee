import "@testing-library/jest-dom/vitest";
import {setupWorker} from "msw/browser"
// import { setupServer } from "msw/node";
import { handlers } from "./mocks/handlers";

import { afterAll, afterEach, beforeAll } from "vitest";

// Start server before all tests

export const server = setupWorker(...handlers);

beforeAll(() => server.start({ onUnhandledRequest: "warn" }));
afterEach(() => server.resetHandlers()); // Reset between tests so individual tests can override
afterAll(() => server.stop());

// export const server = setupServer(...handlers);
// beforeAll(() => server.listen({ onUnhandledRequest: "warn" }));
// afterEach(() => server.resetHandlers()); // Reset between tests so individual tests can override
// afterAll(() => server.close());