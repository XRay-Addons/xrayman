import { type UserStatus, type User } from "@/services/api/generated";

import {
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { enableUserAction, disableUserAction } from "./btn-actions";

import { type VNode } from "vue";

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
    actions.push(
      enableBtn("table.users.actions.enable", enableUserAction(user)),
    );
  }
  if (status !== "disabled") {
    actions.push(
      disableBtn("table.users.actions.disable", disableUserAction(user)),
    );
  }
  actions.push(ensureDeleteBtn("table.users.actions"));

  return mergeActionBtns(actions);
}
