export function pt<T extends Element>(el: ParentNode, name: string): T {
  const pt = el.querySelector(`[data-part="${name}"]`);
  if (!pt) {
    throw new Error(`<${el}> missing data-part="${name}"`);
  }
  return pt as T;
}

export function pts<T extends Element>(el: ParentNode, name: string): T[] {
  const elements = el.querySelectorAll(`[data-part="${name}"]`);
  if (elements.length === 0) {
    throw new Error(`Element missing data-part="${name}"`);
  }
  return Array.from(elements) as T[];
}
