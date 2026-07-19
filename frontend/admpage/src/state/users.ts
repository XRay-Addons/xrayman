import { ref } from "vue";
import type { UserView } from "@/services/api/generated/types.gen";

export const users = ref<UserView[]>([]);
export const usersLoading = ref(false);
