import { User, UserAPIData, UserID } from "./types";
import { config } from "../config/config";

function toAbsoluteUrl(url: string): string {
  return new URL(url, window.location.href).toString();
}

export const ProfileURL = {
  make(user: User): string {
    const path = `${config.USERPAGE_URLPATH}/${user.id}-${user.name}`;
    const absPath = toAbsoluteUrl(path);
    return absPath;
  },
  set(user: User): void {
    history.pushState(null, "", this.make(user));
  },
  reset() {
    const path = `${config.USERPAGE_URLPATH}`;
    const absPath = toAbsoluteUrl(path);
    history.pushState(null, "", absPath);
  },
  parse(): UserID | null {
    const prefix = toAbsoluteUrl(config.USERPAGE_URLPATH);
    const path = window.location.href;
    const match = path.match(new RegExp(`${prefix}/(\\d+)-(.+)$`));
    if (!match) return null;
    const [, idStr, name] = match;
    return { id: Number(idStr), name };
  },
};

const STORAGE_KEY = "user";

export const ProfileStorage = {
  set(user: User): void {
    const data: UserID = { id: user.id, name: user.name };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
  },
  reset() {
    localStorage.removeItem(STORAGE_KEY);
  },
  parse(): UserID | null {
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
  },
};

export const UserUtils = {
  subscriptionURL(user: User): string {
    const subPath = `${config.API_URLPATH}/sub/${user.id}-${user.name}`;
    return toAbsoluteUrl(subPath);
  },

  makeUser(userData: UserAPIData): User {
    return {
      id: userData.id,
      name: userData.name,
      displayName: userData.displayName,
    };
  },
};
