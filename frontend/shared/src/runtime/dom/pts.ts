// methods for safe and type-safe access to element parts.
// use attr data-pt="name" to access via pt<ElemType>(el, name)
// use attr data-pt-private="name"  to access via pt<ElemType>(el, name, "private")

// get element by data-part attribute
export function pt<T extends Element>(el: ParentNode, name: string, suffix: string = ""): T {
  const attr = suffix === "" ? "data-pt" : `data-pt-${suffix}`;
  const pt = el.querySelector(`[${attr}="${name}"]`);
  if (!pt) {
    throw new Error(`<${el}> missing ${attr}="${name}"`);
  }
  return pt as T;
}

// get elements by data-part attribute
export function pts<T extends Element>(el: ParentNode, name: string, suffix: string = ""): T[] {
  const attr = suffix === "" ? "data-pt" : `data-pt-${suffix}`;
  const elements = el.querySelectorAll(`[${attr}="${name}"]`);
  return Array.from(elements) as T[];
}
