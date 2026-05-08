import { ref } from "vue";
import type { Node } from "@/services/api/generated/types.gen";

export const nodes = ref<Node[]>([]);
export const nodesLoading = ref(false);
