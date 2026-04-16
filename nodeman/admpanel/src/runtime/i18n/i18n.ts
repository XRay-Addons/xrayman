import { createI18n } from "vue-i18n";
import en from "@/data/i18n/en.json";
import ru from "@/data/i18n/ru.json";

export const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
  messages: {
    en: en,
    ru: ru,
  },
});

export const t = (text) => {
  return i18n.global.t(text);
};
