import { config } from "../config/config";

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

export const PathTools = {
  admpagePath(path: string): string {
    return makeURL([config.ADMPAGE_URLPATH, path]);
  },
  apiPath(path: string): string {
    return makeURL([config.API_URLPATH, path]);
  },
};
