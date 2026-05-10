<template>
  <ExtendedTable
    :data-source="nodes"
    :loading="nodesLoading"
    :columns="nodesColumns"
    :row-key="rowKey"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import ExtendedTable from "@/vue/components/primitives/table-ext/TableExt.vue";
import { type Node } from "@/services/api/generated/types.gen";
import { useNodesTableColumns } from "./use-columns";
import { onMounted, onBeforeUnmount } from "vue";
import { reloadNodes } from "@/actions/nodes";
import { nodes, nodesLoading } from "@/state/nodes";
import { createPoll } from "@/runtime/polling/server-poll";

// row key
const rowKey = (record: Node): string => String(record.ID);

// i18n prefix
const i18nPrefix = "table.nodes";

// columns
const nodesColumns = useNodesTableColumns(i18nPrefix);

// init data and auto-update
const poll = createPoll(() => reloadNodes({ quiet: true }));
onMounted(() => {
  reloadNodes();
  poll.start();
});
onBeforeUnmount(() => {
  poll.stop();
});
</script>
