export type dtype = "int" | "number" | "string";

export interface DataInputElement extends HTMLElement {
  get value(): number | string;
  set value(v: number | string);
}
