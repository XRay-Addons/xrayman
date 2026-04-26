export type OnClickFn = () => Promise<void>;

export interface ButtonElement extends HTMLElement {
  onClick(fn: OnClickFn): void;
}
