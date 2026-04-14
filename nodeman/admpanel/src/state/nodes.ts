import { ref } from "vue";
import type { Node as APINode } from "../api/generated/types.gen";
import { listNodes } from "../api/client";

// state
export const nodes = ref<APINode[]>([]);
export const nodesLoading = ref(false);
export const nodesError = ref<string | null>(null);

// actions
export async function reloadNodes() {
  nodesLoading.value = true;
  nodesError.value = null;

  try {
    const result = await listNodes();

    if (result.ok) {
      nodes.value = result.data;
    } else {
      nodesError.value = result.reason ?? "Unknown error";
      console.error("Loading nodes error:", result.reason);
    }
  } catch (error) {
    nodesError.value = String(error);
    console.error("Loading nodes error:", error);
  } finally {
    nodesLoading.value = false;
  }
}
