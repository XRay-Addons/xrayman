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
import ExtendedTable, { type ExtendedColumn } from "../ui/TableExt.vue";
import type { User as APIUser } from "../../api/generated/types.gen";
import { listUsers } from "../../api/client";
import { h, ref, onMounted, type VNode } from "vue";
import {
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "../../lib/table-ext-elements";

/* =======================
   state
======================= */

const users = ref<APIUser[]>([]);
const usersLoading = ref(false);

/* =======================
   data loading
======================= */

const loadUsers = async () => {
  usersLoading.value = true;
  try {
    const result = await listUsers();
    if (result.ok) {
      users.value = result.data;
    } else {
      console.error("Loading users error:", result.reason);
    }
  } catch (error) {
    console.error("Loading users error:", error);
  } finally {
    usersLoading.value = false;
  }
};

onMounted(loadUsers);

/* =======================
   row key
======================= */

const rowKey = (record: APIUser): string => String(record.Profile.ID);

/* =======================
   columns (IMPORTANT FIX)
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
    customRender: ({ text }) => renderTag(text),
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
    customRender: ({ text }) =>
      h("span", { style: { fontFamily: "monospace" } }, text),
    extended: true,
  },
  {
    key: "actions",
    dataIndex: ["TargetStatus"],
    customRender: ({ text }) => renderBtns(text),
    extended: true,
  },
];

/* =======================
   helpers
======================= */

function renderTag(text: string) {
  if (text === "enabled") {
    return enabledTag("table.users.status.enabled");
  } else if (text === "disabled") {
    return disabledTag("table.users.status.disabled");
  } else {
    return unknownTag("table.users.status.unknown");
  }
}

function renderBtns(text: string) {
  const actions: VNode[] = [];

  if (text !== "enabled") {
    actions.push(enableBtn("table.users.actions.enable"));
  }
  if (text !== "disabled") {
    actions.push(disableBtn("table.users.actions.disable"));
  }
  actions.push(ensureDeleteBtn("table.users.actions"));

  return mergeActionBtns(actions);
}
</script>
