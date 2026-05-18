import { UserID } from "./user-id";
import { type User } from "@/services/api/generated/types.gen";
import { MakePageUrl } from "@/runtime/utils/paths";

export const ProfileURL = {
  async make(user: User): Promise<string> {
    return MakePageUrl(`${user.Profile.ID}-${user.Profile.Name}`);
  },
  async set(user: User) {
    history.pushState(null, "", await this.make(user));
  },
  async reset() {
    const absPath = await MakePageUrl(`./`);
    history.pushState(null, "", absPath);
  },
  async parse(): Promise<UserID | null> {
    const prefix = await MakePageUrl(`./`);
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
