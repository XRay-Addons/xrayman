import { createI18n } from "vue-i18n";
import en from "@/data/i18n/en.json";
import ru from "@/data/i18n/ru.json";
import { getLanguageState, type Language } from "@/state/language";

const i18n = createI18n({
  legacy: false,
  locale: getLanguageState(),
  fallbackLocale: "en",
  messages: {
    en: en,
    ru: ru,
  },
});

export function updateI18nLanguage() {
  i18n.global.locale.value = getLanguageState();
}

export const t = (text: string): string => {
  return i18n.global.t(text);
};
