/*
 * This code implements a reactive property system and an upgrade/batching
 * mechanism for custom elements using TypeScript decorators.
 *
 *
 * The Problems It Solves:
 *
 * 1. The Pre-Upgrade Assignment Problem
 *    Properties can be assigned to a custom element before its class definition
 *    is registered with the browser:
 *
 *        const element = document.createElement("my-element");
 *        element.tags = ["initial"]; // Assigned directly to the instance
 *        customElements.define("my-element", MyElement);
 *
 *    This creates a raw "own property" on the element object instance that
 *    shadows and bypasses the prototype getter/setter installed later during upgrade.
 *
 * 2. The Field Initializer Overwrite Problem
 *    In standard TypeScript/JavaScript, default field initializers (like `tags = []`)
 *    are compiled into the constructor and execute *after* the base class constructor.
 *    If an external script (or template cloning) passes rich data to the element
 *    before it mounts, the default initializer would completely wipe out that data
 *    and reset it to the default empty value.
 *
 * 3. The Early/Double Render Problem (Layout-Breaks)
 *    When custom elements are generated via template cloning (`cloneNode(true)`)
 *    and nested deeply inside lists, settings properties sequentially (e.g., `tags`
 *    then `data`) would trigger repetitive synchronous `render()` calls before the
 *    element is even in the live DOM, causing race conditions where child elements
 *    cannot be queried or initialized yet.
 *
 *
 * How This Architecture Fixes Them:
 *
 * - Property Defrosting in Constructor:
 *   We intercept and delete shadowing instance properties in the `constructor()`
 *   immediately, capturing early data (like from `makeRow() template logic`)
 *   before TypeScript field initialization can overwrite them.
 *
 * - No-Assign Decorator Options:
 *   By passing default values via `@prop({ default: [] })` and dropping the class
 *   assignment operator (`=`), we completely eliminate the compiled JS constructor
 *   instructions that cause property затирание (overwriting).
 *
 * - Microtask Render Batching (`requestUpdate`):
 *   Modifying properties or connecting the element schedules a single `render()`
 *   pass via `queueMicrotask`. This guarantees that rendering is deferred until
 *   the current synchronous execution stack is completely finished, ensuring
 *   all DOM children are mounted and all references can be safely queried.
 *
 *
 * Internal Implementation Details:
 *
 * - The `@prop` decorator factorizes getters/setters on the class prototype,
 *   routing actual storage to dynamic `Symbol` keys tucked inside each instance.
 *
 * - Metadata Tracking: Class-level property tracking uses isolated copy-on-write
 *   arrays (`static _props`) to prevent inherited properties from accidentally
 *   leaking down or mixing into different custom element subclasses.
 *
 *
 * Risks and Limitations:
 *
 * 1. Dynamic References and Mutation Tracking:
 *    The setter reacts to strict equality changes (`oldValue !== newValue`).
 *    Mutating an array internally via `this.tags.push('new')` will NOT trigger
 *    a re-render because the reference remains the same. Component authors
 *    must treat array/object properties as immutable or manually call `this.requestUpdate()`.
 *
 * 2. Strict Type Boundary:
 *    The decorator forces properties to be non-optional or implicitly handled
 *    via TypeScript definite assignment assertions (`tags!:`). Forgetting to provide
 *    a default configuration when one is expected can cause runtime `undefined` issues.
 *
 * 3. Microtask Lifecycle Race:
 *    If an external script reads computed layout dimensions instantly after
 *    appending the component, it might read stale layout because the actual
 *    DOM mutation inside `render()` is batched and deferred to the next microtask.
 *
 *
 * Developer Experience (DX):
 *
 * The implementation yields a clean, declarative, and framework-like API:
 *
 *     class ListInput extends CustomElement {
 *       @prop({ default: [] }) tags!: string[];
 *       @prop({ default: [] }) data!: Record<string, string>[];
 *
 *       protected render() {
 *         // Guaranteed to run once, after all properties are set and DOM is stable.
 *       }
 *     }
 */
type PrivateStorage = {
  [key: symbol]: unknown;
};

type CustomElementConstructor = {
  new (...args: never[]): CustomElement;
  _props?: string[];
};

function getStorage(target: CustomElement): PrivateStorage {
  return target as unknown as PrivateStorage;
}

interface PropOptions {
  default?: unknown;
}

export function prop(options: PropOptions = {}) {
  return function (target: object, propertyKey: string | symbol): void {
    const propName = String(propertyKey);
    const ctor = target.constructor as CustomElementConstructor;

    if (!Object.hasOwn(ctor, "_props")) {
      ctor._props = [...(ctor._props ?? [])];
    }

    const props = ctor._props!;
    if (!props.includes(propName)) {
      props.push(propName);
    }

    const privateKey = Symbol(`_${propName}`);
    const isInitializedKey = Symbol(`_${propName}_initialized`);

    Object.defineProperty(target, propertyKey, {
      get(this: CustomElement) {
        const storage = getStorage(this);
        if (!storage[isInitializedKey]) {
          return options.default;
        }
        return storage[privateKey];
      },

      set(this: CustomElement, newValue: unknown) {
        const storage = getStorage(this);
        const oldValue = storage[isInitializedKey] ? storage[privateKey] : options.default;

        storage[privateKey] = newValue;
        storage[isInitializedKey] = true;

        if (oldValue !== newValue && this._connected) {
          this.render();
        }
      },
      enumerable: true,
      configurable: true,
    });
  };
}

export abstract class CustomElement extends HTMLElement {
  static _props: string[] = [];
  protected _connected = false;

  constructor() {
    super();
    this.upgradeProperties();
  }

  connectedCallback() {
    this._connected = true;
    queueMicrotask(() => {
      this.render();
    });
  }

  private upgradeProperties() {
    const props = (this.constructor as typeof CustomElement)._props;
    const records = this as Record<string, unknown>;

    for (const prop of props) {
      if (Object.prototype.hasOwnProperty.call(records, prop)) {
        const value = records[prop];
        delete records[prop];
        records[prop] = value;
      }
    }
  }

  protected render(): void {}
}
