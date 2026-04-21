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
import { onMounted } from "vue";
import { reloadNodes } from "@/actions/nodes";
import { nodes, nodesLoading } from "@/state/nodes";

onMounted(reloadNodes);

// row key
const rowKey = (record: Node): string => String(record.ID);

// i18n prefix
const i18nPrefix = "table.nodes";

// columns
const nodesColumns = useNodesTableColumns(i18nPrefix);
</script>
