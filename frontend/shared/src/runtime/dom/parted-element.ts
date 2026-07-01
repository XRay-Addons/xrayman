export abstract class PartedElement extends HTMLElement {
  protected ptEl<T extends Element>(name: string): T {
    return ptEl<T>(this, name);
  }
  protected async ptWebComp<T extends Element>(name: string): Promise<T> {
    return ptWebComp<T>(this, name);
  }
}

export function ptEl<T extends Element>(el: ParentNode, name: string): T {
  const pt = el.querySelector(`[data-part="${name}"]`);
  if (!pt) {
    throw new Error(`<${el}> missing data-part="${name}"`);
  }
  //const tagName = pt.tagName.toLowerCase();
  //if (tagName.includes("-")) {
  //  throw new Error(`<${el}> is custom element, use ptWebComp`);
  //}
  return pt as T;
}

export async function ptWebComp<T extends Element>(el: ParentNode, name: string): Promise<T> {
  const pt = el.querySelector(`[data-part="${name}"]`);
  if (!pt) {
    throw new Error(`<${el}> missing data-part="${name}"`);
  }
  // wait for custom element defined
  const tagName = pt.tagName.toLowerCase();
  if (tagName.includes("-")) {
    await customElements.whenDefined(tagName);
  } else {
    throw new Error(`<${el}> is not custom element, use ptEl`);
  }
  // ensure element defined
  //if (!pt.matches(":defined")) {
  //  customElements.upgrade(pt);
  //}

  return pt as T;
}
