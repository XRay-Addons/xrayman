export type ConfigData = Record<string, string | Record<string, string>[]>;

export interface ConfigEditorElement extends HTMLElement {
  set tags(val: string[]);
  set data(val: ConfigData);
  get value(): ConfigData;
}
