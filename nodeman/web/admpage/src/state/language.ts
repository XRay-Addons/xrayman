export type Language = "en" | "ru";

let language: Language = "ru";

export function getLanguageState() {
  return language;
}

export function setLanguageState(l: Language) {
  language = l;
}
