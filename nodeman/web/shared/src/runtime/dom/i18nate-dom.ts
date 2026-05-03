import { i18nateTree, type T } from "./i18nate-tools";

export function i18nateDOM(t: T) {
  i18nateTree(document.documentElement, t);
}
