import en from "@/data/i18n/en.json";
import ru from "@/data/i18n/ru.json";

export type Language = "en" | "ru";

let language: Language = "ru";

export function setLanguageState(l: Language) {
  language = l;
}

type Messages = {
  [key: string]: string | string[] | Messages;
};

const messages: Messages = { ru, en } as const;

const getValue = (text: string): unknown => {
  return text.split(".").reduce<unknown>((obj, key) => {
    if (typeof obj !== "object" || obj === null) {
      return undefined;
    }

    return (obj as Record<string, unknown>)[key];
  }, messages[language]);
};

export const t = (text: string): string => {
  const value = getValue(text);
  return typeof value === "string" ? value : text;
};

export const tlist = (text: string): string[] => {
  const value = getValue(text);
  return Array.isArray(value) && value.every((item) => typeof item === "string") ? value : [text];
};
