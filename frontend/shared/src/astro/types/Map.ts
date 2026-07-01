export type Item = {
  Key: string;
  Value: string;
};

export interface MapElement extends HTMLElement {
  setValue(val: Item[]): void;
}
