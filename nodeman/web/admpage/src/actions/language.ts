import { setLanguageState, type Language } from "@/state/language";
import { updateI18nLanguage } from "@/runtime/i18n";
import { i18nateDOM } from "@/runtime/dom/i18nate-dom";

export function setLanguage(l: Language) {
  setLanguageState(l);
  updateI18nLanguage();
  i18nateDOM();
}
