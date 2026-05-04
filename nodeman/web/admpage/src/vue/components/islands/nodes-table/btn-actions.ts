import { startNode, stopNode, deleteNode } from "@/services/api/client";
import { reloadNodes } from "@/actions/nodes";
import { type Node } from "@/services/api/generated/types.gen";
import { type BtnAction } from "../../primitives/table-ext/render-primitives";
import { notifyApiError } from "@/runtime/notifications/use-notifications";

export function startNodeAction(node: Node): BtnAction {
  return async () => {
    const r = await startNode(node.ID);
    if (r.ok) {
      reloadNodes({ quiet: true });
    } else {
      notifyApiError("start_node", r.reason);
    }
  };
}

export function stopNodeAction(node: Node): BtnAction {
  return async () => {
    const r = await stopNode(node.ID);
    if (r.ok) {
      reloadNodes({ quiet: true });
    } else {
      notifyApiError("stop_node", r.reason);
    }
  };
}

export function deleteNodeAction(node: Node): BtnAction {
  return async () => {
    const r = await deleteNode(node.ID);
    if (r.ok) {
      reloadNodes({ quiet: true });
    } else {
      notifyApiError("delete_node", r.reason);
    }
  };
}
