import { defineConfig } from "astro/config";
import purgecss from "astro-purgecss";
import vue from "@astrojs/vue";
import relativeLinks from "astro-relative-links";
import compress from "astro-compress";
import swc from "vite-plugin-swc-transform";
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
      JavaScript: false,
      CSS: true,
      HTML: true,
    }),
    vue({
      appEntrypoint: "./src/vue/entry.ts",
    }),
  ],

  vite: {
    plugins: [
      swc({
        include: [/\.ts$/, /astro&type=script/],
        swcOptions: {
          // 1. Force SWC to output modern JavaScript
          jsc: {
            target: "es2022",
            parser: {
              syntax: "typescript",
              decorators: true,
            },
            transform: {
              legacyDecorator: true,
              useDefineForClassFields: false,
            },
          },
        },
      }),
    ],

    ssr: {
      noExternal: ["@xrayman/shared"],
    },

    resolve: {
      alias: {
        "@": path.resolve("./src"),
        "@xrayman/shared": path.resolve("../shared/src"),
      },
    },

    build: {
      target: "es2022",
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
