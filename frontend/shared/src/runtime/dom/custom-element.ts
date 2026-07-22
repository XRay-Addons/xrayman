/*
 * This code implements a small property upgrade mechanism for custom elements.
 *
 * The problem it solves:
 *
 * Custom elements have a specific lifecycle defined by the browser. A property
 * can be assigned to an element before the custom element class is upgraded:
 *
 *     const element = document.createElement("my-element");
 *     element.value = "initial value";
 *     customElements.define("my-element", MyElement);
 *
 * At the moment the assignment happens, `element.value` is just a normal
 * property stored directly on the HTMLElement instance. It does not go through
 * the class setter defined by the custom element.
 *
 * Later, when the browser upgrades the element, the class prototype with
 * accessors/setters becomes active, but the already assigned own property hides
 * the prototype accessor:
 *
 *     instance:
 *       value = "initial value"   <-- own property, shadows the setter
 *
 *     prototype:
 *       set value(...) {}        <-- never called
 *
 * The standard solution is "property upgrade":
 *
 *     1. Remember which properties are managed by the component.
 *     2. When the element is connected, check if those properties were assigned
 *        before upgrade.
 *     3. Remove the own property.
 *     4. Assign the value again, so the real class setter receives it.
 *
 * This allows users of the component to write normal JavaScript:
 *
 *     element.value = "hello";
 *
 * regardless of whether the assignment happened before or after the custom
 * element was upgraded.
 *
 *
 * Why does this implementation look somewhat dirty:
 *
 * - JavaScript decorators are executed while the class is being defined, while
 *   the actual upgrade logic runs on concrete instances later.
 *
 * - The decorator API does not provide a direct communication channel from a
 *   class member decorator to the instance lifecycle methods like
 *   connectedCallback().
 *
 * - Therefore the decorator stores metadata (the list of decorated properties)
 *   on the constructor (`static _props`), and the base class reads this
 *   metadata later.
 *
 * - TypeScript's decorator typings are intentionally generic. The compiler
 *   cannot know that a decorator will only ever be used on classes derived from
 *   CustomElement, so some casts are necessary at this boundary.
 *
 * This is not a runtime type-safe operation in the TypeScript sense; it is a
 * controlled bridge between the decorator system and the browser custom element
 * lifecycle.
 *
 *
 * How other libraries solve this:
 *
 * Many custom element libraries hide this mechanism internally.
 *
 * Examples:
 *
 * - Lit:
 *   LitElement has its own reactive property system. The `@property` decorator
 *   registers metadata during class creation, and the base class performs the
 *   necessary property handling automatically.
 *
 * - FAST:
 *   Microsoft FAST uses observable properties and generated metadata to avoid
 *   manual property management.
 *
 * - Stencil:
 *   Stencil compiles components ahead of time and generates the required
 *   runtime code, avoiding much of the decorator/runtime interaction.
 *
 * - Vanilla Web Components:
 *   The browser specification recommends manually implementing an
 *   `upgradeProperty()` helper in every component constructor.
 *
 * This implementation is basically a small internal version of that pattern.
 *
 *
 * Risks and limitations:
 *
 * 1. Static metadata inheritance:
 *
 *    Static fields are inherited by subclasses. Without creating a separate
 *    `_props` array for each subclass, properties from different components
 *    could accidentally mix.
 *
 * 2. Decorator execution model:
 *
 *    `addInitializer()` runs for every instance, not once when the class is
 *    declared. For a large number of components this means unnecessary work.
 *
 *    A more advanced implementation would use decorator metadata or a separate
 *    class-level registration mechanism.
 *
 * 3. Private fields:
 *
 *    JavaScript private fields (`#field`) are initialized during construction
 *    and cannot be accessed before initialization is complete. Any attempt to
 *    call a setter from an early decorator initializer can fail because the
 *    private storage does not exist yet.
 *
 *    This is why property upgrading is performed from connectedCallback(),
 *    where the instance is already fully initialized.
 *
 * 4. Dynamic property names:
 *
 *    The mechanism relies on runtime property names. Renaming a decorated
 *    accessor without updating metadata generation can silently break the
 *    upgrade process.
 *
 * 5. Encapsulation:
 *
 *    The `_props` field is intentionally public static metadata. A more
 *    encapsulated implementation would use Symbol keys, WeakMap storage, or
 *    the standard decorator metadata proposal.
 *
 *
 * Despite these compromises, this approach provides a simple and predictable
 * developer experience:
 *
 *     class MyInput extends CustomElement {
 *       @prop
 *       accessor value = "";
 *
 *       connectedCallback() {
 *         super.connectedCallback();
 *       }
 *     }
 *
 * The component author only declares reactive properties and uses the normal
 * custom element lifecycle, while all property-upgrade complexity stays inside
 * the base class.
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

export function prop(target: object, propertyKey: string | symbol): void {
  const propName = String(propertyKey);
  const ctor = target.constructor as CustomElementConstructor;

  // Keep separate property metadata for every custom element class.
  // Copy inherited metadata before modifying it to avoid sharing between subclasses.
  if (!Object.hasOwn(ctor, "_props")) {
    ctor._props = [...(ctor._props ?? [])];
  }

  const props = ctor._props!;

  if (!props.includes(propName)) {
    props.push(propName);
  }

  // Store the actual value privately on each instance.
  const privateKey = Symbol(`_${propName}`);

  Object.defineProperty(target, propertyKey, {
    get(this: CustomElement) {
      return getStorage(this)[privateKey];
    },

    set(this: CustomElement, newValue: unknown) {
      const storage = getStorage(this);

      const oldValue = storage[privateKey];
      storage[privateKey] = newValue;

      if (oldValue !== newValue && this._connected) {
        this.render();
      }
    },

    enumerable: true,
    configurable: true,
  });
}

export abstract class CustomElement extends HTMLElement {
  static _props: string[] = [];
  protected _connected = false;

  connectedCallback() {
    this.upgradeProperties();
    queueMicrotask(() => {
      if (!this._connected) {
        this._connected = true;
        this.render();
      }
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
