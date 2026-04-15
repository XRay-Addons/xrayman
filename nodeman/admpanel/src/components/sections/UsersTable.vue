<template>
  <ExtendedTable
    :data-source="users"
    :columns="userColumns"
    :row-key="rowKey"
    :loading="usersLoading"
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
import { onMounted, type VNode, computed } from "vue";
import {
  i18nateColumns,
  makeMonospace,
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
  type BtnAction,
} from "@/lib/table-ext-elements";
import { users, usersLoading, reloadUsers } from "@/state/users";
import { enableUser, disableUser } from "@/api/client";

onMounted(reloadUsers);

// row key
const rowKey = (record: APIUser): string => String(record.Profile.ID);

// i18n prefix
const i18nPrefix = "table.users";

// columns
const userColumns = computed(() => {
  const columns: ExtendedColumn<APIUser>[] = [
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
      customRender: ({ value, record }) => renderActions(value, record),
      extended: true,
    },
  ];
  return i18nateColumns<APIUser>(`${i18nPrefix}.columns`, columns);
});

// value rendering
function renderTag(status: APIUserStatus) {
  if (status === "enabled") {
    return enabledTag("table.users.status.enabled");
  } else if (status === "disabled") {
    return disabledTag("table.users.status.disabled");
  } else {
    return unknownTag("table.users.status.unknown");
  }
}

function enableUserFn(user: APIUser): BtnAction {
  return async () => {
    const r = await enableUser(user.Profile.ID);
    if (r.ok) {
      reloadUsers();
    } else {
      console.log(r.reason);
    }
  };
}

function disableUserFn(user: APIUser): BtnAction {
  return async () => {
    const r = await disableUser(user.Profile.ID);
    if (r.ok) {
      reloadUsers();
    } else {
      console.log(r.reason);
    }
  };
}

function renderActions(status: APIUserStatus, user: APIUser) {
  const actions: VNode[] = [];

  if (status !== "enabled") {
    actions.push(enableBtn("table.users.actions.enable", enableUserFn(user)));
  }
  if (status !== "disabled") {
    actions.push(
      disableBtn("table.users.actions.disable", disableUserFn(user)),
    );
  }
  actions.push(ensureDeleteBtn("table.users.actions"));

  return mergeActionBtns(actions);
}
</script>
