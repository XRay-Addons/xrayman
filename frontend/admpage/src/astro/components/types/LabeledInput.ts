export interface LabeledInputElement extends HTMLElement {
  getValue(): string;
  setValue(val: string): void;
  setPlaceholders(tags: string[]): void;
}
