import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom.js";
import { I18nObserver } from "@xrayman/shared/runtime/dom/i18nate-observer";
import { t } from "@/runtime/i18n";
import { fixIOSActive } from "@xrayman/shared/runtime/dom/ios-active-fix";

i18nateDOM(t);
const i18nObserver = new I18nObserver(t);
i18nObserver.start();
fixIOSActive();

window.addEventListener("beforeunload", () => {
  i18nObserver.stop();
});
