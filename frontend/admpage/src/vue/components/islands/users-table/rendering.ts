import type { UserStatus, User, UserProfile, TrafficStats } from "@/services/api/generated";
import { MakeApiUrl, MakeUserpageURL } from "@/runtime/utils/paths";

import {
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { enableUserAction, disableUserAction, deleteUserAction } from "./btn-actions";

import { makeConfigLine } from "@/vue/components/primitives/table-ext/render-primitives";

import { h, type VNode } from "vue";

export function renderUserPageURL(u: UserProfile) {
  return makeConfigLine(MakeUserpageURL(u.ID, u.Name), true);
}

export function renderTag(status: UserStatus) {
  if (status === "enabled") {
    return enabledTag("table.users.status.enabled");
  } else if (status === "disabled") {
    return disabledTag("table.users.status.disabled");
  } else {
    return unknownTag("table.users.status.unknown");
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

export function renderActions(status: UserStatus, user: User) {
  const actions: VNode[] = [];

  if (status !== "enabled") {
    actions.push(enableBtn("table.users.actions.enable", enableUserAction(user)));
  }
  if (status !== "disabled") {
    actions.push(disableBtn("table.users.actions.disable", disableUserAction(user)));
  }
  actions.push(ensureDeleteBtn("table.users.actions", deleteUserAction(user)));

  return mergeActionBtns(actions);
}

export function renderApiUrl(text: string) {
  text = MakeApiUrl(text);
  return makeConfigLine(text, true);
}
