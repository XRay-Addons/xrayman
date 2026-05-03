import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom.js";
import { I18nObserver } from "@xrayman/shared/runtime/dom/i18nate-observer";
import { setErrorHandler } from "./notifications/errors";
import { notyfErrorHandler } from "./notifications/notyf-handler";
import { t } from "@/runtime/i18n";

setErrorHandler(notyfErrorHandler);

i18nateDOM(t);
const i18nObserver = new I18nObserver(t);
i18nObserver.start();

window.addEventListener("beforeunload", () => {
  i18nObserver.stop();
  setErrorHandler(null);
});
