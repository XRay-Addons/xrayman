import type { ListItemField } from "@xrayman/shared/astro/primitives/ListInput/ListInput.astro";
import type { InputField, ListField } from "@xrayman/shared/astro/primitives/Config/Config.astro";
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
// field names
export const SubscrTitle = "subscr-title";
export const UsersMessage = "users-message";
export const UserPage = "user-page";
export const TgPage = "tg-page";
export const UpdateInterval = "update-interval";
export const AppLinks = "app-links";
export const Headers = "headers";

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
    name: SubscrTitle,
    type: "tagged-input",
    inputProps: { required: false },
  },
  {
    name: UsersMessage,
    type: "tagged-input",
    inputProps: { required: false },
  },
  {
    name: UserPage,
    type: "tagged-input",
    inputProps: { required: false },
  },
  {
    name: TgPage,
    type: "input",
    inputProps: { required: false },
  },
  {
    name: UpdateInterval,
    type: "input",
    inputProps: { required: true, ...intInput },
  },
  {
    name: AppLinks,
    type: "list-input",
    listInputs: appLinkFields,
  },
  {
    name: Headers,
    type: "list-input",
    listInputs: headersFields,
  },
];
