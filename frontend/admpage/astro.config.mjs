import { defineConfig } from "astro/config";
import purgecss from "astro-purgecss";
import vue from "@astrojs/vue";
import relativeLinks from "astro-relative-links";
import compress from "astro-compress";
import path from "node:path";

export default defineConfig({
  output: "static",
  build: {
    outDir: "./dist",
    minify: true,
  },
  integrations: [
    relativeLinks(),
    purgecss(),
    compress({
      JavaScript: true,
      CSS: true,
      HTML: true,
    }),
    vue({
      appEntrypoint: "./src/vue/entry.ts",
    }),
  ],
  vite: {
    resolve: {
      alias: {
        "@": path.resolve("./src"),
        "@xrayman/shared": path.resolve("../shared/src"),
      },
    },
    esbuild: {
      target: "es2022",
    },
    build: {
      target: "es2022",
    },
    experimental: {
      renderBuiltUrl(filename, { hostType, type, ssr }) {
        // for relative paths inside vite framework files
        return { relative: true };
      },
    },
    server: {
      proxy: {
        "/config.js": {
          target: "http://localhost:1001",
          changeOrigin: true,
          rewrite: () => "/adm/config.js",
        },
        "/api": {
          target: "http://localhost:1001",
          changeOrigin: true,
          secure: false,
        },
      },
    },
  },
});
