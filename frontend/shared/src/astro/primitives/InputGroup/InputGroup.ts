export type SubmitHandler = (v: string) => void | Promise<void>;

export interface InputGroupElement extends HTMLElement {
  set submitfn(fn: SubmitHandler);
}
