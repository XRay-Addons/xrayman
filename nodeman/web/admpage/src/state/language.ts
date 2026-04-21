import { ref } from "vue";

export type Language = "en" | "ru";

let language: Language = "en";

export function getLanguageState() {
  return language;
}

export function setLanguageState(l: Language) {
  language = l;
}
