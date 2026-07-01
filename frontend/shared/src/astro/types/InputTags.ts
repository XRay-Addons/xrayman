export interface InputTagsInit {
  tags: string[];
  input: HTMLInputElement;
}

export interface InputTagsElement extends HTMLElement {
  set init(v: InputTagsInit);
}
