import type { dtype } from "./DataInput";

export interface CardItemProp {
  i18n: string;
  dtype: dtype;
  field: string;
  attrs?: astroHTML.JSX.InputHTMLAttributes;
}

export interface ArrayItemProp {
  i18n: string;
  dtype: dtype;
  field: string;
  attrs?: astroHTML.JSX.InputHTMLAttributes;
}

export interface CardArrayProp {
  i18n: string;
  dtype: dtype;
  field: string;
  props: ArrayItemProp[];
  attrs?: astroHTML.JSX.InputHTMLAttributes;
}

export type CardItem = CardItemProp | CardArrayProp;

export interface Column {
  dtype: dtype;
  field: string;
  attrs?: astroHTML.JSX.InputHTMLAttributes;
}
