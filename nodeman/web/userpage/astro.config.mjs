import { defineConfig } from "astro/config";
import { purgecss } from "@zokki/astro-purgecss";
import relativeLinks from "astro-relative-links";
import process from "process";
import compress from "astro-compress";
import path from "node:path";

export default defineConfig({
  output: "static",
  integrations: [
    relativeLinks(),
    purgecss({
      rejected: true,
    }),
    compress({
      JavaScript: true,
      CSS: true,
      HTML: true,
    }),
  ],
  build: {
    outDir: "./dist",
    minify: true,
  },
  vite: {
    resolve: {
      alias: {
        "@": path.resolve("./src"),
        "@xrayman/shared": path.resolve("../shared/src"),
      },
    },
  },
});
