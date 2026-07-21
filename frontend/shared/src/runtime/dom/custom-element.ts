export function prop(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any,
  context: ClassAccessorDecoratorContext,
) {
  const fieldName = context.name;

  context.addInitializer(function () {
    const self = this as Record<PropertyKey, unknown>;

    if (Object.prototype.hasOwnProperty.call(self, fieldName)) {
      const oldValue = self[fieldName];
      delete self[fieldName];
      value.set.call(self, oldValue);
    }
  });

  return {
    set(this: CustomElement, newValue: unknown) {
      const oldValue = value.get.call(this);
      value.set.call(this, newValue);

      if (oldValue !== newValue && this._connected) {
        this.render();
      }
    },
  };
}

export abstract class CustomElement extends HTMLElement {
  protected _connected: boolean = false;

  connectedCallback() {
    queueMicrotask(() => {
      if (!this._connected) {
        this._connected = true;
        this.render();
      }
    });
  }
  protected render(): void {}
}
