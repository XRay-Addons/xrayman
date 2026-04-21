import { enableUser, disableUser } from "@/services/api/client";
import { reloadUsers } from "@/actions/users";
import { type User } from "@/services/api/generated/types.gen";
import { type BtnAction } from "../../primitives/table-ext/render-primitives";
import { notifyApiError } from "@/runtime/notifications/errors";

export function enableUserAction(user: User): BtnAction {
  return async () => {
    const r = await enableUser(user.Profile.ID);
    if (r.ok) {
      reloadUsers({ quiet: true });
    } else {
      notifyApiError("enable_user", r.reason);
    }
  };
}

export function disableUserAction(user: User): BtnAction {
  return async () => {
    const r = await disableUser(user.Profile.ID);
    if (r.ok) {
      reloadUsers({ quiet: true });
    } else {
      notifyApiError("disable_user", r.reason);
    }
  };
}
