import { I18nObserver } from "./i18nate-observer";
import { i18nateTree } from "./i18nate-tools";

const i18nObserver = new I18nObserver();

export function i18nateDOM() {
  i18nateTree(document.documentElement);
}

export function startI18nObserver() {
  i18nObserver.start();
}

export function stopI18nObserver() {
  i18nObserver.stop();
}
