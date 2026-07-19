import en from "@/data/i18n/en.json";
import ru from "@/data/i18n/ru.json";

export type Language = "en" | "ru";

let language: Language = "ru";

export function setLanguageState(l: Language) {
  language = l;
}

type Messages = {
  [key: string]: string | Messages;
};

const messages: Messages = { ru, en } as const;

export const t = (text: string): string => {
  const value = text.split(".").reduce<unknown>((obj, key) => {
    if (typeof obj !== "object" || obj === null) {
      return undefined;
    }

    return (obj as Record<string, unknown>)[key];
  }, messages[language]);

  return typeof value === "string" ? value : text;
};
