import { defineConfig } from "astro/config";
import relativeLinks from "astro-relative-links";
import process from "process";
import vue from "@astrojs/vue";

export default defineConfig({
  output: "static",
  integrations: [
    relativeLinks(),
    vue({
      appEntrypoint: "/src/vue/entry.ts",
    }),
  ],
  outDir: process.env.NODEMAN_WEB_DIST || "./dist",
});
