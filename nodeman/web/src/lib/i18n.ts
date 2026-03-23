import ru from "../data/i18n/ru.json";
import en from "../data/i18n/en.json";

export const messages = { ru, en } as const;
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

function applyI18n(root: ParentNode = document) {
  root.querySelectorAll<HTMLElement>("*").forEach((el) => {
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
  });
}

function startI18nObserver() {
  const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      for (const node of mutation.addedNodes) {
        if (!(node instanceof HTMLElement)) continue;
        applyI18n(node);
      }
    }
  });

  observer.observe(document.body, {
    childList: true,
    subtree: true,
  });
}

function setLang(lang: Lang) {
  currentLang = lang;
  applyI18n();
}

setLang("ru");
startI18nObserver();
