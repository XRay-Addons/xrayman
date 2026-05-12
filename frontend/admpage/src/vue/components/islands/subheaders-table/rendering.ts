import { type Header } from "@/services/api/generated";
import { MakeApiUrl } from "@/runtime/utils/paths";

import {
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { deleteHeaderAction } from "./btn-actions";

import {} from "@/vue/components/primitives/table-ext/render-primitives";

import { type VNode } from "vue";

export function renderActions(header: Header) {
  const actions: VNode[] = [];

  actions.push(ensureDeleteBtn("table.sub-headers.actions", deleteHeaderAction(header)));

  return mergeActionBtns(actions);
}
