import { defineConfig } from "astro/config";
import vue from "@astrojs/vue";

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
        "@": "/src",
      },
    },
  },
});
