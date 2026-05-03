import { defineConfig } from "astro/config";
import vue from "@astrojs/vue";
import path from "node:path";

export default defineConfig({
  output: "static",
  build: {
    outDir: "dist",
  },
  integrations: [
    vue({
      appEntrypoint: "/src/vue/entry.ts",
    }),
  ],
  vite: {
    resolve: {
      alias: {
        "@": path.resolve("./src"),
        "@xrayman/shared": path.resolve("../shared/src"),
      },
    },
  },
});
