import { User, UserAPIData, UserID } from "./types";
import { config } from "../config/config";

export const UserUtils = {
  profilePath(user: User): string {
    return `${config.USERPAGE_URLPATH}/${user.id}-${user.name}`;
  },

  profileURL(user: User): string {
    return `${window.location.origin}${UserUtils.profilePath(user)}`;
  },

  parseProfileURL(): UserID | null {
    const path = window.location.pathname;
    const match = path.match(
      new RegExp(`${config.USERPAGE_URLPATH}/(\\d+)-(.+)$`),
    );
    if (!match) return null;
    const [, idStr, name] = match;
    return { id: Number(idStr), name };
  },

  subscriptionPath(user: User): string {
    return `${config.API_URLPATH}/sub/${user.id}-${user.name}`;
  },

  subscriptionURL(user: User): string {
    return `${window.location.origin}${UserUtils.subscriptionPath(user)}`;
  },

  makeUser(userData: UserAPIData): User {
    return {
      id: userData.id,
      name: userData.name,
      displayName: userData.displayName,
    };
  },
};
