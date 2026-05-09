import { i18nateElement, i18nateTree, type T } from "./i18nate-tools";

export class I18nObserver {
  private observer: MutationObserver | null = null;

  constructor(private readonly t: T) {}

  start() {
    if (this.observer) return;

    this.observer = new MutationObserver((mutations) => {
      mutations.forEach(this.processMutation);
    });

    this.observer.observe(document.documentElement, {
      childList: true,
      subtree: true,
      attributes: false,
    });
  }

  private processMutation = (mutation: MutationRecord) => {
    const { type, addedNodes, target, attributeName } = mutation;

    switch (type) {
      case "childList":
        addedNodes.forEach(
          (node) => node instanceof HTMLElement && i18nateTree(node as HTMLElement, this.t),
        );
        break;
      case "attributes":
        // not actually used for now
        //if (target instanceof HTMLElement) {
        //  i18nateElement(target, this.t);
        //}
        break;
    }
  };

  stop() {
    this.observer?.disconnect();
    this.observer = null;
  }
}
