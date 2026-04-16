import { ref } from "vue";
import type { Node as APINode } from "../api/generated/types.gen";
import { listNodes } from "../api/client";
import { serverErrorNotification } from "@/runtime/notifications/errors";

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
      nodesError.value = result.reason;
      serverErrorNotification("get_nodes", result.reason);
    }
  } finally {
    nodesLoading.value = false;
  }
}

export async function syncNodes() {
  const result = await listNodes();
  if (result.ok) {
    nodes.value = result.data;
  } else {
    console.error("Loading nodes error:", result.reason);
  }
}
