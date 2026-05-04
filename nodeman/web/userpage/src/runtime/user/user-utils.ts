import { UserID } from "./user-id";
import { type User } from "@/services/api/generated/types.gen";
import { MakePageUrl } from "@xrayman/shared/runtime/paths/paths";

export const ProfileURL = {
  make(user: User): string {
    return MakePageUrl(`./${user.Profile.ID}-${user.Profile.Name}`);
  },
  set(user: User): void {
    history.pushState(null, "", this.make(user));
  },
  reset() {
    const absPath = MakePageUrl(`./`);
    history.pushState(null, "", absPath);
  },
  parse(): UserID | null {
    const prefix = MakePageUrl(`./`);
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
    const data: UserID = { id: user.Profile.ID, name: user.Profile.Name };
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
