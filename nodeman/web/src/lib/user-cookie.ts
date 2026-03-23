import { UserID } from "./types";

const STORAGE_KEY = "user";

export function parseCookie(): UserID | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;

    const data = JSON.parse(raw) as UserID;

    if (typeof data.id !== "number" || typeof data.name !== "string") {
      return null;
    }

    return data;
  } catch {
    return null;
  }
}

export function setCookie(id: number, name: string) {
  const data: UserID = { id, name };
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
}

export function resetCookie() {
  localStorage.removeItem(STORAGE_KEY);
}
