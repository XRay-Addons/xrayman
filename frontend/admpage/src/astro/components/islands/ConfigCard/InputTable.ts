import type { dtype } from "./DataInput";

export interface Column {
  dtype: dtype;
  field: string;
  placeholders: boolean;
  attrs?: astroHTML.JSX.InputHTMLAttributes;
}

export interface InputTableElement extends HTMLElement {
  get value(): Record<string, any>;
  set value(v: Record<string, any>);
}
