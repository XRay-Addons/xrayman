import { ref } from "vue";
import type { User } from "../api/generated/types.gen";
import { listUsers } from "../api/client";
import { errorNotification } from "@/runtime/notifications/errors";
import { serverErrorNotification } from "@/runtime/notifications/errors";

// state
export const users = ref<User[]>([]);
export const usersLoading = ref(false);

// actions
export async function reloadUsers() {
  usersLoading.value = true;

  try {
    const result = await listUsers();
    if (result.ok) {
      users.value = result.data;
    } else {
      serverErrorNotification("get_users", result.reason);
      console.error("Loading users error:", result.reason);
    }
  } finally {
    usersLoading.value = false;
  }
}

export async function syncUsers() {
  const result = await listUsers();
  if (result.ok) {
    users.value = result.data;
  } else {
    console.error("Loading users error:", result.reason);
  }
}
