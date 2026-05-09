export type OnClickFn = (val: string) => Promise<void>;

export interface InputGroupElement extends HTMLElement {
  getValue(): string;
  setValue(val: string): void;
  onClick(fn: OnClickFn): void;
  getButton(): HTMLButtonElement;
  getInput(): HTMLInputElement;
}
