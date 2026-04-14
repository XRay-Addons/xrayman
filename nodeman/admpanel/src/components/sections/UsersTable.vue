<template>
  <ExtendedTable
    :data-source="users"
    :columns="userColumns"
    :row-key="rowKey"
    :loading="usersLoading"
    i18n-prefix="table.users"
    color="#ff0000b5"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import ExtendedTable, {
  type ExtendedColumn,
} from "@/components/ui/TableExt.vue";
import type {
  User as APIUser,
  UserStatus as APIUserStatus,
} from "@/api/generated/types.gen";
import { onMounted, type VNode } from "vue";
import {
  makeMonospace,
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/lib/table-ext-elements";
import { users, usersLoading, reloadUsers } from "@/state/users";

onMounted(reloadUsers);

/* =======================
   row key
======================= */

const rowKey = (record: APIUser): string => String(record.Profile.ID);

/* =======================
   columns 
======================= */

const userColumns: ExtendedColumn<APIUser>[] = [
  {
    key: "id",
    dataIndex: ["Profile", "ID"],
    width: "8ch",
  },
  {
    key: "display-name",
    dataIndex: ["Profile", "DisplayName"],
    width: "20ch",
  },
  {
    key: "target-status",
    dataIndex: ["TargetStatus"],
    customRender: ({ value }) => renderTag(value),
  },
  {
    key: "name",
    dataIndex: ["Profile", "Name"],
    extended: true,
  },
  {
    key: "vless-uuid",
    dataIndex: ["Profile", "VlessUUID"],
    ellipsis: true,
    width: "8ch",
    customRender: ({ text }) => makeMonospace(text),
    extended: true,
  },
  {
    key: "actions",
    dataIndex: ["TargetStatus"],
    customRender: ({ value }) => renderActions(value),
    extended: true,
  },
];

/* =======================
   helpers
======================= */

function renderTag(status: APIUserStatus) {
  if (status === "enabled") {
    return enabledTag("table.users.status.enabled");
  } else if (status === "disabled") {
    return disabledTag("table.users.status.disabled");
  } else {
    return unknownTag("table.users.status.unknown");
  }
}

function renderActions(status: APIUserStatus) {
  const actions: VNode[] = [];

  if (status !== "enabled") {
    actions.push(enableBtn("table.users.actions.enable"));
  }
  if (status !== "disabled") {
    actions.push(disableBtn("table.users.actions.disable"));
  }
  actions.push(ensureDeleteBtn("table.users.actions"));

  return mergeActionBtns(actions);
}
</script>
