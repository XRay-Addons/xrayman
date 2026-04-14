<template>
  <ExtendedTable
    :data-source="users"
    :columns="userColumns"
    :rowKey="rowKey"
    i18n-prefix="table.users"
    color="#ff0000b5"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import ExtendedTable from "./ExtendedTable.vue";
import { h, ref, onMounted, onBeforeUnmount } from "vue";
import { Tag, Button, Space, Popconfirm } from "ant-design-vue";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons-vue";

interface User {
  Profile: Profile;
  TargetStatus: "enabled" | "disabled" | "unknown";
}

interface Profile {
  ID: number;
  Name: string;
  DisplayName: string;
  VlessUUID: string;
}

const data: User[] = [
  {
    Profile: {
      ID: 1,
      Name: "alice",
      DisplayName: "Alice Johnson",
      VlessUUID: "123e4567-e89b-12d3-a456-426614174000",
    },
    TargetStatus: "enabled",
  },
  {
    Profile: {
      ID: 2,
      Name: "bob",
      DisplayName: "Bob Smith",
      VlessUUID: "223e4567-e89b-12d3-a456-426614174001",
    },
    TargetStatus: "disabled",
  },
  {
    Profile: {
      ID: 3,
      Name: "charlie",
      DisplayName: "Charlie Brown",
      VlessUUID: "323e4567-e89b-12d3-a456-426614174002",
    },
    TargetStatus: "unknown",
  },
];

const users = ref(data);

const rowKey = (record) => record.Profile.ID;

const userColumns = [
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

    customRender: ({ text }) => {
      return renderTag(text);
    },
  },
  {
    key: "name",
    dataIndex: ["Profile", "Name"],
    extended: true,
  },
  {
    key: "vless-uuid",
    dataIndex: ["Profile", "VlessUUID"],
    extended: true,
    ellipsis: true,
    width: "8ch",
  },
  {
    key: "actions",
    dataIndex: ["TargetStatus"],
    customRender: ({ text }) => {
      return renderBtns(text);
    },
    extended: true,
  },
];

function renderTag(text) {
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

function makeTag(color: string, i18n: string, icon) {
  return h(
    Tag,
    {
      color: color,
    },
    {
      default: () => [h(icon), h("span", { "data-i18n": i18n })],
    },
  );
}

function renderBtns(text) {
  const actions = [];
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
      title: h("span", {
        "data-i18n": "table.users.actions.delete-confirn-header",
      }),
      description: h("span", {
        "data-i18n": "table.users.actions.delete-confirn-body",
      }),
      okText: h("span", {
        "data-i18n": "table.users.actions.delete-confirn-yes",
      }),
      cancelText: h("span", {
        "data-i18n": "table.users.actions.delete-confirn-no",
      }),
    },
    {
      default: () => {
        return h(Button, {
          danger: true,
          size: "small",
          type: "primary",
          style: { boxShadow: "none" },
          "data-i18n": "table.users.actions.delete",
        });
      },
    },
  );
}
</script>
