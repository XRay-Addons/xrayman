const I18N_PREFIX = "data-i18n";

export type T = (s: string) => string;

export function i18nateElement(el: HTMLElement, t: T) {
  for (const attr of Array.from(el.attributes)) {
    if (!attr.name.startsWith(I18N_PREFIX)) continue;
    // translate value
    const value = t(attr.value as string);
    // attr may point to direct value (if name == prefix)
    // or to another attribute:
    // data-i18n-placeholder -> placeholder
    // data-i18n-data-key -> data-key
    if (attr.name == I18N_PREFIX) {
      el.innerHTML = value;
    } else {
      el.setAttribute(attr.name.replace(`${I18N_PREFIX}-`, ""), value);
    }
  }
}

export function i18nateTree(root: HTMLElement, t: T) {
  i18nateElement(root, t);
  root.querySelectorAll<HTMLElement>("*").forEach((el) => i18nateElement(el, t));
}
