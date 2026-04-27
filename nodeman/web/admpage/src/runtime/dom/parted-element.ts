export abstract class PartedElement extends HTMLElement {
  protected pt<T extends Element>(name: string): T {
    const el = this.querySelector(`[data-part="${name}"]`);
    if (!el) {
      throw new Error(`<${this.tagName}> missing data-part="${name}"`);
    }
    return el as T;
  }
}
