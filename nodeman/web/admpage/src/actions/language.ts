import { setLanguageState, type Language } from "@/state/language";
import { updateI18nLanguage } from "@/runtime/i18n";
import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom";
import { type T } from "@xrayman/shared/runtime/dom/i18nate-tools";

export function setLanguage(l: Language, t: T) {
  setLanguageState(l);
  updateI18nLanguage();
  i18nateDOM(t);
}
