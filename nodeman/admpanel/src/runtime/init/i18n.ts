import { applyI18nTree, startI18nObserver } from "@/runtime/i18n/dom-observer";

export function initI18nObserver() {
  applyI18nTree(document.root);
  startI18nObserver(document.root);
}
