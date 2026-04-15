import { initI18nObserver } from "./i18n";
import { fixIOSActive } from "./ios-active-fix";

function initRuntime() {
  initI18nObserver();
  fixIOSActive();
}

initRuntime();
