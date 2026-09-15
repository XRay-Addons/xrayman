import { CustomElement } from "@xrayman/shared/runtime/dom/custom-element";

export class DialogItemElement extends CustomElement {
  async onShow(): Promise<void> {}
}

export interface DialogElement extends Element {
  showModal(): void | Promise<void>;
}
