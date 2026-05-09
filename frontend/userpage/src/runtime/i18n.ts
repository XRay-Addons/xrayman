import en from "@/data/i18n/en.json";
import ru from "@/data/i18n/ru.json";

export type Language = "en" | "ru";

let language: Language = "ru";

export function setLanguageState(l: Language) {
  language = l;
}

const messages = { ru, en } as const;

export const t = (text: string): string => {
  return messages[language][text] ?? text;
};
