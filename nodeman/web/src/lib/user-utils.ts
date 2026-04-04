import { User, UserAPIData, UserID } from "./types";
import { PathTools } from "./path-utils";

export const ProfileURL = {
  make(user: User): string {
    return PathTools.userpagePath(`./${user.id}-${user.name}`);
  },
  set(user: User): void {
    history.pushState(null, "", this.make(user));
  },
  reset() {
    const absPath = PathTools.userpagePath(`./`);
    history.pushState(null, "", absPath);
  },
  parse(): UserID | null {
    const prefix = PathTools.userpagePath(`./`);
    const path = window.location.href;
    const match = path.match(new RegExp(`${prefix}(\\d+)-(.+)$`));
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
    return PathTools.apiPath(`./sub/${user.id}-${user.name}`);
  },

  makeUser(userData: UserAPIData): User {
    return {
      id: userData.id,
      name: userData.name,
      displayName: userData.displayName,
    };
  },
};
