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
import ExtendedTable from "../ui/ExtendedTable.vue";
import { h, ref, onMounted } from "vue";
import type { VNode } from "vue";
import { Tag, Button, Space, Popconfirm } from "ant-design-vue";
import type { ColumnType } from "ant-design-vue/es/table";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons-vue";
import type { User as APIUser } from "../../api/generated/types.gen";
import { listUsers } from "../../api/client";
import type { ExtendedColumn } from "../ui/ExtendedTable.vue";

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
    return makeTag(
      "success",
      "table.users.status.enabled",
      CheckCircleOutlined,
    );
  } else if (text === "disabled") {
    return makeTag("error", "table.users.status.disabled", CloseCircleOutlined);
  } else {
    return makeTag(
      "warning",
      "table.users.status.unknown",
      ExclamationCircleOutlined,
    );
  }
}

function makeTag(color: string, i18n: string, icon: any) {
  return h(
    Tag,
    { color },
    {
      default: () => [h(icon), h("span", { "data-i18n": i18n })],
    },
  );
}

function renderBtns(text: string) {
  const actions: VNode[] = [];

  if (text === "enabled") {
    actions.push(disableBtn());
  } else if (text === "disabled") {
    actions.push(enableBtn());
  } else {
    actions.push(enableBtn(), disableBtn());
  }

  actions.push(deleteBtn());

  return h(Space, { size: "small" }, () => actions);
}

function enableBtn() {
  return h(Button, {
    ghost: true,
    size: "small",
    type: "primary",
    "data-i18n": "table.users.actions.enable",
  });
}

function disableBtn() {
  return h(Button, {
    danger: true,
    ghost: true,
    size: "small",
    type: "primary",
    "data-i18n": "table.users.actions.disable",
  });
}

function deleteBtn() {
  return h(
    Popconfirm,
    {
      okText: h("span", {
        "data-i18n": "table.users.actions.delete-confirn-yes",
      }),
      cancelText: h("span", {
        "data-i18n": "table.users.actions.delete-confirn-no",
      }),
    },
    {
      default: () =>
        h(Button, {
          danger: true,
          size: "small",
          type: "primary",
          style: { boxShadow: "none" },
          "data-i18n": "table.users.actions.delete",
        }),
      title: () =>
        h("span", {
          "data-i18n": "table.users.actions.delete-confirn-header",
        }),
      description: () =>
        h("span", {
          "data-i18n": "table.users.actions.delete-confirn-body",
        }),
    },
  );
}
</script>
