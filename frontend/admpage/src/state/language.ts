import type { Language } from "@xrayman/shared/runtime/dom/i18n";

let language: Language = "en";

export function getLanguageState(): Language {
  return language;
}

export function setLanguageState(l: Language) {
  language = l;
}
