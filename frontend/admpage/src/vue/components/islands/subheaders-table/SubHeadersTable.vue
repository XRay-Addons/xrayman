<template>
  <ExtendedTable
    :data-source="subHeaders"
    :loading="subHeadersLoading"
    :columns="subHeadersColumns"
    :row-key="rowKey"
    :scroll="{ x: 'max-content' }"
    :pagination="false"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import ExtendedTable from "@/vue/components/primitives/table-ext/TableExt.vue";
import { type Header } from "@/services/api/generated/types.gen";
import { useSubHeadersTableColumns } from "./use-columns";
import { onMounted, onBeforeUnmount } from "vue";
import { reloadSubHeaders } from "@/actions/sub-headers";
import { subHeaders, subHeadersLoading } from "@/state/sub-headers";
import { createPoll } from "@/runtime/polling/server-poll";

// row key
const rowKey = (record: Header): string => String(record.ID);

// i18n prefix
const i18nPrefix = "table.sub-headers";

// columns
const subHeadersColumns = useSubHeadersTableColumns(i18nPrefix);

// init data and auto-update
const poll = createPoll(() => reloadSubHeaders({ quiet: true }));
onMounted(() => {
  reloadSubHeaders();
  poll.start();
});
onBeforeUnmount(() => {
  poll.stop();
});
</script>
