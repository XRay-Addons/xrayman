import { i18nateElement } from "./i18nate-tools";

export class I18nObserver {
  private observer: MutationObserver | null = null;

  start() {
    if (this.observer) return;

    this.observer = new MutationObserver((mutations) => {
      mutations.forEach(processMutation);
    });

    this.observer.observe(document.documentElement, {
      childList: true,
      subtree: true,
      attributes: false,
    });
  }

  stop() {
    this.observer?.disconnect();
    this.observer = null;
  }
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
