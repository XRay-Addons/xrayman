/**
 * Runtime helper for styles that cannot be expressed dynamically in CSS alone.
 *
 * CSS is declarative, but some layout rules depend on runtime values:
 * viewport size, element measurements, browser state, etc.
 *
 * This class provides a minimal bridge between CSS classes and JavaScript
 * calculations. A handler is registered for a class name and is executed for
 * every matching element whenever the layout may have changed:
 *
 * - window resize
 * - dynamically inserted DOM nodes
 *
 * Usage:
 *
 *   sizeObserver.add("fluid-md", (el) => {
 *     el.style.setProperty("--fluid-width", "50");
 *   });
 *
 * HTML stays purely declarative:
 *
 *   <div class="fluid-md"></div>
 *
 * The goal is not to replace CSS, but to extend it in places where CSS has no
 * way to express a required dynamic calculation.
 */
export type Handler = (el: HTMLElement) => void;

export class SizeObserver {
  private handlers = new Map<string, Set<Handler>>();
  private elements = new Map<string, Set<HTMLElement>>();

  private readonly mutationObserver: MutationObserver;
  private readonly resizeHandler = () => this.update();

  constructor() {
    window.addEventListener("resize", this.resizeHandler);

    this.mutationObserver = new MutationObserver((records) => {
      for (const record of records) {
        for (const node of record.addedNodes) {
          if (node instanceof HTMLElement) {
            this.scan(node);
          }
        }
      }

      this.update();
    });

    this.mutationObserver.observe(document.body, {
      childList: true,
      subtree: true,
    });

    this.scan(document);
  }

  add(className: string, handler: Handler) {
    let handlers = this.handlers.get(className);

    if (!handlers) {
      handlers = new Set();
      this.handlers.set(className, handlers);
    }

    handlers.add(handler);

    this.scan(document);
    this.update();
  }

  stop() {
    window.removeEventListener("resize", this.resizeHandler);
    this.mutationObserver.disconnect();

    this.handlers.clear();
    this.elements.clear();
  }

  private scan(root: ParentNode) {
    for (const className of this.handlers.keys()) {
      const elements = this.elements.get(className) ?? new Set<HTMLElement>();

      if (root instanceof HTMLElement && root.classList.contains(className)) {
        elements.add(root);
      }

      root.querySelectorAll<HTMLElement>(`.${className}`).forEach((el) => elements.add(el));

      this.elements.set(className, elements);
    }
  }

  private update() {
    for (const [className, elements] of this.elements) {
      const handlers = this.handlers.get(className);

      if (!handlers) {
        continue;
      }

      for (const el of elements) {
        for (const handler of handlers) {
          handler(el);
        }
      }
    }
  }
}
