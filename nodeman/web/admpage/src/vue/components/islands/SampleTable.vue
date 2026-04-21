<template>
  <TableExtended :data-source="users" :columns="columns" row-key="id" />
</template>

<script setup lang="ts">
import TableExtended from "@/vue/components/primitives/table-ext/TableExt.vue";
import type { ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import {
  enabledTag,
  makeMonospace,
  makeCopyable,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { computed } from "vue";

// mock data
interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  age: number;
}

const users: User[] = [
  {
    id: "1",
    name: "Alice Johnson",
    email: "alice@example.com",
    role: "admin",
    age: 29,
  },
  {
    id: "2",
    name: "Bob Smith",
    email: "bob@example.com",
    role: "user",
    age: 34,
  },
  {
    id: "3",
    name: "Charlie Brown",
    email: "charlie@example.com",
    role: "moderator",
    age: 22,
  },
];

// column schema
const columns = computed(() => {
  const c: ExtendedColumn<User>[] = [
    {
      key: "name",
      dataIndex: "name",
    },
    {
      key: "email",
      dataIndex: "email",
    },
    {
      key: "role",
      dataIndex: "role",
      extended: true,
      customRender: ({ value }) => enabledTag(value),
    },
    {
      key: "age",
      dataIndex: "age",
      extended: true,
      customRender: ({ value }) => makeCopyable(makeMonospace(value), value),
    },
  ];
  return i18nateColumns<User>(`table.users.columns`, c);
});
</script>
