import ru from "@/data/i18n/ru.json";
//import en from "@/data/i18n/en.json";

export const messages = { ru } as const;
export type Lang = keyof typeof messages;
export type I18nKey = keyof typeof ru;

let currentLang: Lang = "ru";

function t(key: I18nKey): string {
  return messages[currentLang][key] ?? key;
}

const I18N_PREFIX = "data-i18n";
const DATA_PREFIX = "data";

function isSpecialAttr(attr: string) {
  return attr == "placeholder" || attr == "value";
}

function applyI18nEl(el: HTMLElement) {
  for (const attr of Array.from(el.attributes)) {
    if (!attr.name.startsWith(I18N_PREFIX)) continue;

    const key = attr.value as I18nKey;
    const value = t(key);

    if (attr.name == I18N_PREFIX) {
      el.innerHTML = value;
      continue;
    }

    const target = attr.name.replace(`${I18N_PREFIX}-`, "");
    const targetAttr = isSpecialAttr(target)
      ? target
      : `${DATA_PREFIX}-${target}`;
    el.setAttribute(targetAttr, value);
  }
}

function applyI18nTree(root: Node) {
  if (root instanceof HTMLElement) {
    applyI18nEl(root);
    root.querySelectorAll<HTMLElement>("*").forEach(applyI18nEl);
  }
}

export function startI18nObserver() {
  applyI18nTree(document.body);

  const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.type === "childList") {
        for (const node of mutation.addedNodes) {
          applyI18nTree(node);
        }
      }

      if (
        mutation.type === "attributes" &&
        mutation.oldValue == null &&
        mutation.target instanceof HTMLElement
      ) {
        applyI18nEl(mutation.target);
      }
    }
  });

  observer.observe(document.body, {
    childList: true,
    subtree: true,
    attributes: true,
    attributeOldValue: true,
  });
}
