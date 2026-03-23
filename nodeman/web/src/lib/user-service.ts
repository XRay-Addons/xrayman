import { UserAPIData, User } from "./types";
import { userStore } from "./user-store";
import { parseCookie } from "./user-cookie";
import { parseURL, makeURL } from "./user-url";

export function deriveUser(data: UserAPIData): User {
  return {
    id: data.id,
    name: data.name,
    displayName: data.displayName,
    subscriptionURL: data.subscriptionURL,
    userPageURL: makeURL(data.id, data.name),
  };
}

/*export async function initUser() {
  const userID = parseURL() ?? parseCookie();

  try {
    const res = await fetchUser(userID);

    const res = await fetch(`/api/users/${userID}`);
    if (!res.ok) throw new Error("Failed to fetch user");
    res.status

    const apiData: UserAPIData = await res.json();
    const user = deriveUser(apiData);

    userStore.set(user);
  } catch (err) {
    console.error(err);
    userStore.set(null);
  }
}*/
