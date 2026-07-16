// methods for safe and type-safe access to element parts.
// use data-pt attribute for elements of nested components
// use data-ppt (p for private) for elements of current component to avoid conflicts

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

// get element by data-ppart attribute
export function ppt<T extends Element>(el: ParentNode, name: string): T {
  const pt = el.querySelector(`[data-ppt="${name}"]`);
  if (!pt) {
    throw new Error(`<${el}> missing data-ppt="${name}"`);
  }
  return pt as T;
}

// get elements by data-ppart attribute
export function ppts<T extends Element>(el: ParentNode, name: string): T[] {
  const elements = el.querySelectorAll(`[data-ppt="${name}"]`);
  return Array.from(elements) as T[];
}
