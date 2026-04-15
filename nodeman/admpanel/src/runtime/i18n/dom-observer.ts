import { i18n } from "@/runtime/i18n/i18n";

const I18N_PREFIX = "data-i18n";
const DATA_PREFIX = "data";

function t(key: string): string {
  return i18n.global.t(key);
}

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

export function applyI18nTree(root: Node) {
  if (root instanceof HTMLElement) {
    applyI18nEl(root);
    root.querySelectorAll<HTMLElement>("*").forEach(applyI18nEl);
  }
}

export function startI18nObserver() {
  // create ans start DOM changing observer
  const domChangeObserver = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      // process new child
      if (mutation.type === "childList") {
        for (const node of mutation.addedNodes) {
          applyI18nTree(node);
        }
      }

      // process atributes changing
      if (
        mutation.type === "attributes" &&
        mutation.oldValue == null &&
        mutation.target instanceof HTMLElement
      ) {
        applyI18nEl(mutation.target);
      }
    }
  });

  domChangeObserver.observe(document.documentElement, {
    childList: true,
    subtree: true,
    attributes: true,
    attributeOldValue: true,
  });

  // start lang-change event observer
  const langChangeObserver = window.addEventListener(
    "language-changed",
    (event) => {
      console.log("on lang changed");
      applyI18nTree(document.documentElement);
    },
  );
}
