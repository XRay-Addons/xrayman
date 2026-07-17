import type { HTMLAttributes } from "astro/types";
import type { WithPrefix } from "@xrayman/shared/astro/props/props";

// input attributes set for common inputs: integer, string, number
export const strInput = { type: "text" } satisfies HTMLAttributes<"input">;

export const intInput = {
  type: "number",
  inputmode: "numeric",
  step: "1",
  pattern: "[0-9]*",
} satisfies HTMLAttributes<"input">;

export const numInput = {
  type: "number",
  inputmode: "decimal",
  step: "any",
  pattern: "[0-9]*([\\.,][0-9]+)?",
} satisfies HTMLAttributes<"input">;

// add prefix input: to input attrs
export type InputAttrs = HTMLAttributes<"input">;

export function withInputPrefix(attrs: InputAttrs): WithPrefix<InputAttrs, "input:"> {
  return Object.fromEntries(
    Object.entries(attrs).map(([key, value]) => [`input:${key}`, value]),
  ) as WithPrefix<InputAttrs, "input:">;
}
