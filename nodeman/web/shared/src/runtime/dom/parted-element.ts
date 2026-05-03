export abstract class PartedElement extends HTMLElement {
  protected pt<T extends Element>(name: string): T {
    return pt<T>(this, name);
  }
}

export function pt<T extends Element>(el: HTMLElement, name: string): T {
  const pt = el.querySelector(`[data-part="${name}"]`);
  if (!pt) {
    throw new Error(`<${el.tagName}> missing data-part="${name}"`);
  }
  return pt as T;
}
