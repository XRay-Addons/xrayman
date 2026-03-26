import { User, UserAPIData, UserID } from "./types";
import { config } from "../config/config";

function toAbsoluteUrl(url: string): string {
  return new URL(url, window.location.href).toString();
}

export const UserUtils = {
  profilePath(user: User): string {
    return `${config.USERPAGE_URLPATH}/${user.id}-${user.name}`;
  },

  profileURL(user: User): string {
    return toAbsoluteUrl(UserUtils.profilePath(user));
  },

  parseProfileURL(): UserID | null {
    const prefix = toAbsoluteUrl(config.USERPAGE_URLPATH);
    const path = window.location.href;
    const match = path.match(new RegExp(`${prefix}/(\\d+)-(.+)$`));
    if (!match) return null;
    const [, idStr, name] = match;
    return { id: Number(idStr), name };
  },

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
