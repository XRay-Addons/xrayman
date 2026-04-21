import { ref } from "vue";
import type { User } from "@/services/api/generated/types.gen";

export const users = ref<User[]>([]);
export const usersLoading = ref(false);
