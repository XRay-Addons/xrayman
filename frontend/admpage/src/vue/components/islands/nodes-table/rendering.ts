import { type NodeStatus, type Node, type TrafficStats } from "@/services/api/generated";

import {
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { startNodeAction, stopNodeAction, deleteNodeAction } from "./btn-actions";

import { type VNode, h } from "vue";

export function renderTag(status: NodeStatus) {
  if (status === "stopped") {
    return disabledTag("table.nodes.status.stopped");
  } else if (status === "running") {
    return enabledTag("table.nodes.status.running");
  } else {
    return unknownTag("table.nodes.status.unknown");
  }
}

export function renderTraffic(traffic: TrafficStats) {
  const total = traffic.Download + traffic.Upload;
  return h("span", {}, trafficText(total));
}

function trafficText(traffic: number): string {
  if (traffic == 0) {
    return "0";
  }

  const suffixes = ["B", "KB", "MB", "GB", "TB"];
  let suffixIdx = Math.floor(Math.log(traffic) / Math.log(1024));
  suffixIdx = Math.min(suffixIdx, suffixes.length - 1);

  const value = traffic / Math.pow(1024, suffixIdx);
  return `${value.toFixed(1)} ${suffixes[suffixIdx]}`;
}

export function renderLoad(load: number) {
  const value = Math.round(load);
  return h("span", {}, `${value}%`);
}

export function renderActions(status: NodeStatus, node: Node) {
  const actions: VNode[] = [];

  if (status !== "running") {
    actions.push(enableBtn("table.nodes.actions.start", startNodeAction(node)));
  }
  if (status !== "stopped") {
    actions.push(disableBtn("table.nodes.actions.stop", stopNodeAction(node)));
  }
  actions.push(ensureDeleteBtn("table.nodes.actions", deleteNodeAction(node)));

  return mergeActionBtns(actions);
}
