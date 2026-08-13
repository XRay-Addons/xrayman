import { defineConfig } from "astro/config";
import purgecss from "astro-purgecss";
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
      JavaScript: true,
      CSS: true,
      HTML: true,
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
    //experimental: {
    //  renderBuiltUrl(filename, { hostType, type, ssr }) {
    //    // for relative paths inside vite framework files
    //    return { relative: true };
    //  },
    //},
    // use local go server for backend
    server: {
      proxy: {
        "/config.js": {
          target: "http://localhost:1001",
          changeOrigin: true,
          rewrite: () => "/u/config.js",
        },
        "/u/config.js": {
          target: "http://localhost:1001",
          changeOrigin: true,
          rewrite: () => "/u/config.js",
        },
        "/api": {
          target: "http://localhost:1001",
          changeOrigin: true,
          secure: false,
        },
        "/u/": {
          target: "http://localhost:4321",
          changeOrigin: true,
          rewrite: () => "",
        },
      },
    },
  },
  server: {
    port: 4321,
    strictPort: true,
  },
});
