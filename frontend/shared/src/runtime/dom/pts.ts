// methods for safe and type-safe access to element parts.
// use attr data-pt="name" to access via pt<ElemType>(el, name)
// use attr data-item-name  to access via item<ElemType>(el, "item-name")

// get element by data-part attribute
export function pt<T extends Element>(el: ParentNode, name: string): T {
  const pt = el.querySelector(`[data-pt="${name}"]`);
  if (!pt) {
    throw new Error(`<${el}> missing data-pt="${name}"`);
  }
  return pt as T;
}

// get elements by data-part attribute
export function pts<T extends Element>(el: ParentNode, name: string): T[] {
  const elements = el.querySelectorAll(`[data-pt="${name}"]`);
  return Array.from(elements) as T[];
}

export function item<T extends Element>(el: ParentNode, name: string): T {
  const item = el.querySelector(`[data-${name}], [data-pt="${name}"]`);
  if (!item) {
    throw new Error(`<${el}> missing data-${name} or data-pt="${name}"`);
  }
  return item as T;
}

export function items<T extends Element>(el: ParentNode, name: string): T[] {
  const items = el.querySelectorAll(`[data-${name}], [data-pt="${name}"]`);
  return Array.from(items) as T[];
}
