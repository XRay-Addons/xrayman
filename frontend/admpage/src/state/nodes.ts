import { ref } from "vue";
import type { NodeView } from "@/services/api/generated/types.gen";

export const nodes = ref<NodeView[]>([]);
export const nodesLoading = ref(false);
