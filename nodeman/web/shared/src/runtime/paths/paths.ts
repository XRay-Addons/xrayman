//  API_URLPATH: (window as any).__API_URL__ ?? "/api",
//  USERPAGE_URLPATH: (window as any).__SPA_URL__ ?? "/",

function makeURL(components: string[]): string {
  // remove start and end "/"
  for (var i = 0; i < components.length; i++) {
    components[i] = components[i].replace(/^\/|\/$/g, "");
  }
  return new URL(components.join("/"), window.location.origin).toString();
}

export function MakeFullUrl(prefix: string, path: string): string {
  return makeURL([prefix, path]);
}
