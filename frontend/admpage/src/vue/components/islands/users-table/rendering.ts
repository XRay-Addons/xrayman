import { type UserStatus, type User, UserProfile } from "@/services/api/generated";
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

import { type VNode } from "vue";

export async function renderUserPageURL(u: UserProfile) {
  return makeConfigLine(await MakeUserpageURL(u.ID, u.Name), true);
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

export async function renderApiUrl(text: string) {
  text = await MakeApiUrl(text);
  return makeConfigLine(text, true);
}
