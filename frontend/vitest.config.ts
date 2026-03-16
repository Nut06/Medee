import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import path from "path";
import { playwright } from "@vitest/browser-playwright";

export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [["babel-plugin-react-compiler"]],
      },
    }),
  ],

  test: {
    browser: {
      provider: playwright(),
      enabled: true,
      headless: true,
      instances: [{ browser: "chromium" }],
    },
    globals: true,
    // environment: "jsdom",
    setupFiles: "./src/test/setup.ts",
    css: true,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
