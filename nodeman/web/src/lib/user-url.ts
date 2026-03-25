import { UserID } from "./types";

export function parseURL(): UserID | null {
  const path = window.location.pathname;
  const match = path.match(/^\/(\d+)-(.+)$/);
  if (!match) return null;

  const [, idStr, name] = match;
  return { id: Number(idStr), name };
}

export function setURL(url: string) {
  history.pushState(null, "", `/${url}`);
}

export function resetURL() {
  history.pushState(null, "", "");
}
