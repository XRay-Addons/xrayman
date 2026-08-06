import { i18nateDOM } from "@xrayman/shared/runtime/dom/i18nate-dom.js";
import { I18nObserver } from "@xrayman/shared/runtime/dom/i18nate-observer";
import { t } from "@/runtime/i18n";
import { SizeObserver } from "@xrayman/shared/runtime/dom/size-observer";
import { fixIOSActive } from "@xrayman/shared/runtime/dom/ios-active-fix";
import { initTStyles } from "@xrayman/shared/styles/tstyles";

i18nateDOM(t);
const i18nObserver = new I18nObserver(t);
i18nObserver.start();
fixIOSActive();

const sizeObserver = new SizeObserver();
initTStyles(sizeObserver);

window.addEventListener("beforeunload", () => {
  i18nObserver.stop();
  sizeObserver.stop();
});
