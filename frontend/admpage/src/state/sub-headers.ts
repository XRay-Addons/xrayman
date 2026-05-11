import { ref } from "vue";
import type { Header } from "@/services/api/generated/types.gen";

export const subHeaders = ref<Header[]>([]);
export const subHeadersLoading = ref(false);
