/*import { User, UserAPIData, UserID } from "./types";
import { PathTools } from "./path-utils";

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
};*/
