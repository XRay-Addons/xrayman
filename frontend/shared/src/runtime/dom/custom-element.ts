export function prop(value: any, context: ClassAccessorDecoratorContext) {
  const fieldName = context.name;

  context.addInitializer(function (this: any) {
    if (Object.prototype.hasOwnProperty.call(this, fieldName)) {
      const oldValue = this[fieldName];
      delete this[fieldName];
      value.set.call(this, oldValue);
    }
  });

  return {
    set(this: any, newValue: any) {
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
