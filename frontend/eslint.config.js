import eslintPluginAstro from "eslint-plugin-astro";
import tseslint from "typescript-eslint";

export default tseslint.config(
  {
    ignores: [
      "**/dist/**",
      "**/.astro/**",
      "**/node_modules/**",
      "**/*.config.*",
      "**/tsconfig*.json",
      "**/src/services/api/generated/**",
    ],
  },

  ...tseslint.configs.recommended,
  ...eslintPluginAstro.configs.recommended,

  {
    files: ["**/*.astro"],
    rules: {
      "astro/no-unused-define-vars-in-style": "error",
      "astro/valid-compile": "error",
    },
  },
);
