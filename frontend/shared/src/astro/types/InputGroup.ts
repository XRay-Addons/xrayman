export type OnClickFn = (val: string) => Promise<void>;

export interface InputGroupElement extends HTMLElement {
  get value(): string;
  set value(val: string);
  set onclickfn(fn: OnClickFn);
  get button(): HTMLButtonElement;
  get input(): HTMLInputElement;
}
