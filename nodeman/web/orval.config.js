module.exports = {
  api: {
    input: "../../xrayman/nodeman/pkg/api/http/openapi/openapi.yaml",
    output: {
      target: "./src/api/generated.ts",
      client: "axios",
      mode: "single",
    },
  },
};
