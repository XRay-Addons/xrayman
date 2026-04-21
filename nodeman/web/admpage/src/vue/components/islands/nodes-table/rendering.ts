import { type NodeStatus, type Node } from "@/services/api/generated";

import {
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { startNodeAction, stopNodeAction } from "./btn-actions";

import { type VNode } from "vue";

export function renderTag(status: NodeStatus) {
  if (status === "stopped") {
    return disabledTag("table.nodes.status.stopped");
  } else if (status === "running") {
    return enabledTag("table.nodes.status.running");
  } else {
    return unknownTag("table.nodes.status.unknown");
  }
}

export function renderActions(status: NodeStatus, node: Node) {
  const actions: VNode[] = [];

  if (status !== "running") {
    actions.push(enableBtn("table.nodes.actions.start", startNodeAction(node)));
  }
  if (status !== "stopped") {
    actions.push(disableBtn("table.nodes.actions.stop", stopNodeAction(node)));
  }
  actions.push(ensureDeleteBtn("table.nodes.actions"));

  return mergeActionBtns(actions);
}
