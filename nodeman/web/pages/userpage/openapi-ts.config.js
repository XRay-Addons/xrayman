import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "../../../pkg/api/http/openapi/openapi.yaml",
  output: "./src/services/api/generated",
  plugins: ["@hey-api/typescript", "@hey-api/sdk", "@hey-api/client-fetch"],
  parser: {
    filters: {
      tags: {
        include: ["userpage"],
      },
    },
  },
});
