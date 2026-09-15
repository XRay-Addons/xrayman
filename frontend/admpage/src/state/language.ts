export type Language = "en";

let language: Language = "en";

export function getLanguageState() {
  return language;
}

export function setLanguageState(l: Language) {
  language = l;
}
