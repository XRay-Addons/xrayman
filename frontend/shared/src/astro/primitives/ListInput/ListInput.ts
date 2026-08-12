import type { HTMLAttributes, ComponentProps } from "astro/types";
import TextInput from "@xrayman/shared/astro/primitives/TextInput.astro";

export type FieldType = "input" | "tagged-input";
export type InputAttrs = Omit<ComponentProps<typeof TextInput>, "name">;

export interface ListItemField {
  name: string;
  type: FieldType;
  inputProps?: InputAttrs;
  width?: string;
}

export interface ListInputProps extends HTMLAttributes<"div"> {
  name: string;
  i18n: string;
  fields: ListItemField[];
}

export interface ListInputElement extends HTMLElement {
  set tags(val: string[]);
  set data(val: Record<string, string>[]);
  get value(): Record<string, string>[];
}
