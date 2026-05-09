/*import { setLanguageState, type Language } from "@/state/language";
import { updateI18nLanguage } from "@/runtime/i18n";
import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom";
import { type T } from "@xrayman/shared/runtime/dom/i18nate-tools";*/

import { type Language, setLanguageState } from "@/runtime/i18n";
import { type T } from "@xrayman/shared/runtime/dom/i18nate-tools";
import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom";

export function setLanguage(l: Language, t: T) {
  setLanguageState(l);
  i18nateDOM(t);
}
