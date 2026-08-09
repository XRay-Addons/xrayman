import type { ListItemField } from "@xrayman/shared/astro/primitives/ListInput/ListInput.astro";
import type { TextInput } from "@xrayman/shared/astro/primitives/TextInput.astro";
import type { InputField, ListField } from "@xrayman/shared/astro/primitives/Config/Config.astro";
import type { ComponentProps } from "astro/types";

// -------------------------------------------------------------------------
// input attributes set for common inputs: integer, string, number
export const strInput = {
  type: "text",
} satisfies ComponentProps<typeof TextInput>;

export const intInput = {
  type: "number",
  inputmode: "numeric",
  step: "1",
  pattern: "[0-9]*",
} satisfies ComponentProps<typeof TextInput>;

export const numInput = {
  type: "number",
  inputmode: "decimal",
  step: "any",
  pattern: "[0-9]*([\\.,][0-9]+)?",
} satisfies ComponentProps<typeof TextInput>;

// -------------------------------------------------------------------------
// field names
export const SubscrTitle = "subscr-title";
export const UsersMessage = "users-message";
export const UserPage = "user-page";
export const TgPage = "tg-page";
export const UpdateInterval = "update-interval";
export const AppLinks = "app-links";
export const Routing = "routing";
export const CustomHeaders = "custom-headers";

// -------------------------------------------------------------------------
const platformAppFields: ListItemField[] = [
  {
    name: "Name",
    type: "input",
    inputProps: { required: true },
    width: "20%",
  },
  {
    name: "Platforms",
    type: "input",
    inputProps: { required: true },
    width: "20%",
  },
  {
    name: "URL",
    type: "input",
    inputProps: { required: true, nowrap: false },
    width: "60%",
  },
];

export const headersFields: ListItemField[] = [
  {
    name: "Key",
    type: "input",
    inputProps: { required: true },
  },
  {
    name: "Value",
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
    inputProps: { required: false, nowrap: false },
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
    name: Routing,
    type: "input",
    inputProps: { required: false, nowrap: false },
  },
  {
    name: AppLinks,
    type: "list-input",
    listInputs: platformAppFields,
  },
  {
    name: CustomHeaders,
    type: "list-input",
    listInputs: headersFields,
  },
];
