function makeURL(components: string[]): string {
  const filtered: string[] = [];
  for (const component of components) {
    filtered.push(component.replace(/^\/+|\/+$/g, ""));
  }
  return new URL(filtered.join("/"), window.location.origin).toString();
}

export function MakeFullUrl(prefix: string, path: string): string {
  return makeURL([prefix, path]);
}
