// methods for safe and type-safe access to element parts.
// use attr data-pt="name" or data-name to access via pt<ElemType>(el, name)

export function pt<T extends Element>(el: ParentNode, name: string): T {
  const item = el.querySelector(`[data-${name}], [data-pt="${name}"]`);
  if (!item) {
    throw new Error(`<${el}> missing data-${name} or data-pt="${name}"`);
  }
  return item as T;
}

export function pts<T extends Element>(el: ParentNode, name: string): T[] {
  const items = el.querySelectorAll(`[data-${name}], [data-pt="${name}"]`);
  return Array.from(items) as T[];
}
