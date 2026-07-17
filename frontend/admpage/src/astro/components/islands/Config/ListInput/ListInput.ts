export interface ListInputElement extends HTMLElement {
  set tags(val: string[]);
  set data(val: Record<string, string>[]);
  get value(): Record<string, string>[];
}
