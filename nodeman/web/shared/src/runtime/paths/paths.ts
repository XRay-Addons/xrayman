//  API_URLPATH: (window as any).__API_URL__ ?? "/api",
//  USERPAGE_URLPATH: (window as any).__SPA_URL__ ?? "/",

function makeURL(components: string[]): string {
  var current = window.location.origin;
  for (var i = 0; i < components.length; i++) {
    if (current.endsWith("/")) {
      current = new URL(components[i], current).toString();
    } else {
      current = new URL(components[i], current + "/").toString();
    }
  }
  return current;
}

export function MakeApiUrl(path: string): string {
  return makeURL([(window as any).__API_URL__ ?? "/api", path]);
}

export function MakePageUrl(path: string): string {
  return makeURL([(window as any).__SPA_URL__ ?? "/u", path]);
}
