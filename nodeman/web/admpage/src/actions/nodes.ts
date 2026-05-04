import { listNodes } from "@/services/api/client";
import { nodes, nodesLoading } from "@/state/nodes";
import { notifyApiError } from "@/runtime/notifications/use-notifications";

export async function reloadNodes({
  quiet = false,
}: {
  quiet?: boolean;
} = {}) {
  if (quiet) {
    const result = await listNodes();
    if (result.ok) nodes.value = result.data;
    return;
  }

  nodesLoading.value = true;

  try {
    const result = await listNodes();

    if (result.ok) {
      nodes.value = result.data;
    } else {
      notifyApiError("get_nodes", result.reason);
    }
  } finally {
    nodesLoading.value = false;
  }
}
