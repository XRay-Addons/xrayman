import { t } from "@/runtime/i18n";

const I18N_PREFIX = "data-i18n";

function i18nateElement(el: HTMLElement) {
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

function i18nateTree(root: HTMLElement) {
  i18nateElement(root);
  root.querySelectorAll<HTMLElement>("*").forEach(i18nateElement);
}

function processMutation({ type, addedNodes, target }: MutationRecord) {
  switch (type) {
    case "childList":
      addedNodes.forEach(
        (node) =>
          node instanceof HTMLElement && i18nateElement(node as HTMLElement),
      );
      break;
    case "attributes":
      if (target instanceof HTMLElement) {
        i18nateElement(target);
      }
      break;
  }
}

class I18nObserver {
  private observer: MutationObserver | null = null;

  start() {
    if (this.observer) return;

    this.observer = new MutationObserver((mutations) => {
      mutations.forEach(processMutation);
    });

    this.observer.observe(document.documentElement, {
      childList: true,
      subtree: true,
      attributes: true,
    });
  }

  stop() {
    this.observer?.disconnect();
    this.observer = null;
  }
}

const i18nObserver = new I18nObserver();

export function startI18nObserver() {
  i18nObserver.start();
}

export function stopI18nObserver() {
  i18nObserver.stop();
}

export function updateDOMLanguage() {
  i18nateTree(document.documentElement);
}
