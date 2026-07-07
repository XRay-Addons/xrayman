import type { dtype } from "./DataInput";

export interface ArrayItemProp {
  i18n: string;
  dtype: dtype;
  field: string;
  attrs?: astroHTML.JSX.InputHTMLAttributes;
}
