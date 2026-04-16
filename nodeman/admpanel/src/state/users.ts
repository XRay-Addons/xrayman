import { ref } from "vue";
import type { User as APIUser } from "../api/generated/types.gen";
import { listUsers } from "../api/client";

// state
export const users = ref<APIUser[]>([]);
export const usersLoading = ref(false);
export const usersError = ref<string | null>(null);

// actions
export async function reloadUsers() {
  usersLoading.value = true;
  usersError.value = null;

  const result = await listUsers();
  if (result.ok) {
    users.value = result.data;
  } else {
    usersError.value = result.reason ?? "Unknown error";
    console.error("Loading users error:", result.reason);
  }

  usersLoading.value = false;
}

export async function syncUsers() {
  const result = await listUsers();
  if (result.ok) {
    users.value = result.data;
  } else {
    console.error("Loading users error:", result.reason);
  }
}
