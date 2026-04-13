<template>
  <ExtendedTable
    :data-source="users"
    :columns="userColumns"
    :rowKey="rowKey"
    i18n-prefix="table.users"
  />
</template>

<script setup lang="ts">
import ExtendedTable from "./ExtendedTable.vue";
import { h, ref, onMounted, onBeforeUnmount } from "vue";
import { Tag, Button, Space } from "ant-design-vue";

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
    ellipsis: true,
    width: "8ch",
  },
  {
    key: "display-name",
    dataIndex: ["Profile", "DisplayName"],
    ellipsis: true,
  },
  {
    key: "target-status",
    dataIndex: ["TargetStatus"],
    ellipsis: true,
    customRender: ({ text }) => {
      return renderTag(text);
    },
  },
  {
    key: "name",
    dataIndex: ["Profile", "Name"],
    extended: true,
    ellipsis: true,
  },
  {
    key: "vless-uuid",
    dataIndex: ["Profile", "VlessUUID"],
    extended: true,
    ellipsis: true,
  },
  {
    key: "actions",
    dataIndex: ["TargetStatus"],
    customRender: ({ text }) => {
      return renderBtns(text);
    },
    width: "20ch",
    ellipsis: true,
    extended: true,
  },
];

function renderTag(text) {
  if (text === "enabled") {
    return enabledTag();
  } else if (text === "disabled") {
    return disabledTag();
  } else {
    return unknownTag();
  }
}

function enabledTag() {
  return h(Tag, {
    color: "success",
    "data-i18n": "table.users.status.enabled",
  });
}

function disabledTag() {
  return h(Tag, {
    color: "error",
    "data-i18n": "table.users.status.disabled",
  });
}

function unknownTag() {
  return h(Tag, {
    color: "default",
    "data-i18n": "table.users.status.unknown",
  });
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
    type: "primary",
    ghost: true,
    size: "small",
    "data-i18n": "table.users.actions.enable",
  });
}

function disableBtn() {
  return h(Button, {
    type: "primary",
    ghost: true,
    danger: true,
    size: "small",
    "data-i18n": "table.users.actions.disable",
  });
}

function deleteBtn() {
  return h(Button, {
    type: "primary",
    danger: true,
    size: "small",
    "data-i18n": "table.users.actions.delete",
  });
}
</script>
