<template>
  <a-config-provider :theme="theme">
    <a-table
      class="table-ext"
      :data-source="dataSource"
      :columns="mainCols"
      :row-key="rowKey"
      :expanded-row-render="expandedRowRender"
      size="medium"
      v-bind="$attrs"
    >
    </a-table>
  </a-config-provider>
</template>

<script setup lang="ts" generic="T">
import { getPaletteRef } from "@/state/palette";
import { useTableTheme } from "./use-table-theme";
import { useTableColumns } from "./use-table-columns";
import { useExpandedRowRender } from "./use-expanded-render";
import { type Props } from "./table-types";

// props
const props = defineProps<Props<T>>();

// theme
const theme = useTableTheme(getPaletteRef());

// columns
const { mainCols, extCols } = useTableColumns(props.columns);

// expand
const keyColumnWidth = props.expandedKeyColumnWidth ?? "16ch";
const expandedRowRender = useExpandedRowRender(extCols.value, keyColumnWidth);
</script>

<style scoped src="./table-style.scss" lang="scss"></style>
