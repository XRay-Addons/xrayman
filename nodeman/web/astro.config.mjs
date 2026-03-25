import { defineConfig } from "astro/config";
import relativeLinks from "astro-relative-links";
import process from "process";

export default defineConfig({
  output: "static",
  integrations: [relativeLinks()],
  outDir: process.env.NODEMAN_WEB_DIST || "./dist",
});
