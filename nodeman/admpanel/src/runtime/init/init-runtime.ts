import { startI18nObserver } from "./i18n";
import { fixIOSActive } from "./ios-active-fix";

function initRuntime() {
  startI18nObserver();
  fixIOSActive();
}

initRuntime();
