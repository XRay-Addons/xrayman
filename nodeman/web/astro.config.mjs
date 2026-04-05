import { defineConfig } from "astro/config";
import { purgecss } from "@zokki/astro-purgecss";
import relativeLinks from "astro-relative-links";
import process from "process";

export default defineConfig({
  output: "static",
  integrations: [relativeLinks(), purgecss()],
  outDir: process.env.NODEMAN_WEB_DIST || "./dist",
});
