<template>
  <ExtendedTable
    :data-source="users"
    :loading="usersLoading"
    :columns="usersColumns"
    :row-key="rowKey"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import ExtendedTable from "@/vue/components/primitives/table-ext/TableExt.vue";
import { type User } from "@/services/api/generated/types.gen";
import { useUsersTableColumns } from "./use-columns";
import { onMounted } from "vue";
import { reloadUsers } from "@/actions/users";
import { users, usersLoading } from "@/state/users";

onMounted(reloadUsers);

// row key
const rowKey = (record: User): string => String(record.Profile.ID);

// i18n prefix
const i18nPrefix = "table.users";

// columns
const usersColumns = useUsersTableColumns(i18nPrefix);
</script>
