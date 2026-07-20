export type SubmitHandler = (v: string) => void | Promise<void>;

export interface InputGroupElement extends HTMLElement {
  set submitfn(fn: SubmitHandler);
}

/*import { pt } from "@xrayman/shared/runtime/dom/pts";

export function onSubmit(el: HTMLFormElement, fn: (value: string) => void | Promise<void>) {
  const inputEl = pt<HTMLInputElement>(el, "input-group-input");
  const btnEl = pt<HTMLButtonElement>(el, "input-group-btn");

  btnEl.addEventListener("click", (e: Event) => {
    e.preventDefault();
    el.requestSubmit();
  });
  el.addEventListener("submit", async (e: Event) => {
    e.preventDefault();
    await fn(inputEl.value);
  });
}*/
