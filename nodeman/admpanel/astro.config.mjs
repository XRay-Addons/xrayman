import { defineConfig } from "astro/config";
import { purgecss } from "@zokki/astro-purgecss";
import relativeLinks from "astro-relative-links";
import process from "process";
import compress from "astro-compress";
import vue from "@astrojs/vue";
import path from "node:path";

export default defineConfig({
  output: "static",
  integrations: [
    relativeLinks(),
    purgecss(),
    compress({
      JavaScript: true,
      CSS: true,
      HTML: true,
    }),
    vue({
      appEntrypoint: "/src/vue/entry.ts",
    }),
  ],
  vite: {
    resolve: {
      alias: {
        "@": path.resolve("./src"),
      },
    },
  },
  outDir: process.env.NODEMAN_WEB_DIST || "./dist",
  build: {
    minify: true,
  },
});
