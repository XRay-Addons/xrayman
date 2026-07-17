import type { ListItemField } from "../Config/ListInput/ListInput.astro";
import type { InputField, ListField } from "../Config/Config.astro";
import type { HTMLAttributes } from "astro/types";

// -------------------------------------------------------------------------
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

// -------------------------------------------------------------------------

const appLinkFields: ListItemField[] = [
  {
    name: "description",
    type: "input",
    inputProps: { required: true },
    width: "20%",
  },
  {
    name: "platforms",
    type: "input",
    inputProps: { required: true },
    width: "20%",
  },
  {
    name: "url",
    type: "input",
    inputProps: { required: true },
    width: "60%",
  },
];

export const headersFields: ListItemField[] = [
  {
    name: "key",
    type: "input",
    inputProps: { required: true },
  },
  {
    name: "value",
    type: "tagged-input",
    inputProps: { required: false },
  },
];

export const fields: (InputField | ListField)[] = [
  {
    name: "subscr-title",
    type: "tagged-input",
    inputProps: { required: false },
  },
  {
    name: "users-message",
    type: "tagged-input",
    inputProps: { required: false },
  },
  {
    name: "user-page",
    type: "tagged-input",
    inputProps: { required: false },
  },
  {
    name: "tg-page",
    type: "input",
    inputProps: { required: false },
  },
  {
    name: "update-interval",
    type: "input",
    inputProps: { required: true, ...intInput },
  },
  {
    name: "app-links",
    type: "list-input",
    listInputs: appLinkFields,
  },
  {
    name: "headers",
    type: "list-input",
    listInputs: headersFields,
  },
];
