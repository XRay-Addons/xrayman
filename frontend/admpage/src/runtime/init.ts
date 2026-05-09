import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom.js";
import { I18nObserver } from "@xrayman/shared/runtime/dom/i18nate-observer";
import { t } from "@/runtime/i18n";

i18nateDOM(t);
const i18nObserver = new I18nObserver(t);
i18nObserver.start();

window.addEventListener("beforeunload", () => {
  i18nObserver.stop();
});
