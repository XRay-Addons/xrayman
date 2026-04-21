import { listUsers } from "@/services/api/client";
import { users, usersLoading } from "@/state/users";
import { notifyApiError } from "@/runtime/notifications/errors";

export async function reloadUsers({
  quiet = false,
}: {
  quiet?: boolean;
} = {}) {
  if (quiet) {
    const result = await listUsers();
    if (result.ok) users.value = result.data;
    return;
  }

  usersLoading.value = true;

  try {
    const result = await listUsers();

    if (result.ok) {
      users.value = result.data;
    } else {
      notifyApiError("get_users", result.reason);
    }
  } finally {
    usersLoading.value = false;
  }
}
