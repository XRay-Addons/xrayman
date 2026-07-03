export type InputType = "text" | "integer";

export type RWMapData = Record<string, string>;

export interface RWMapElement extends HTMLElement {
  set value(val: RWMapData);
  set placeholders(val: string[]);
  getValue(): RWMapData;
}
