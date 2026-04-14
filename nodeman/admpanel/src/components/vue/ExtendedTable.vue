<template>
  <a-config-provider :theme="theme">
    <a-table
      :data-source="dataSource"
      :columns="mainColumns"
      :row-key="rowKey"
      :scroll="{ x: 'max-width' }"
      size="medium"
      class="table-extended"
      v-bind="$attrs"
    >
      <!-- Header -->
      <template #headerCell="{ column }">
        <span v-if="column.key" :data-i18n="`${i18nPrefix}.${column.key}`">
          {{ column.key }}
        </span>
      </template>

      <!-- Expanded row -->
      <template #expandedRowRender="{ record }">
        <a-table
          :columns="extendedTableColumns"
          :data-source="extendedDataSource(record)"
          :pagination="false"
          :show-header="false"
          size="small"
          class="table-extension"
        />
      </template>
    </a-table>
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed, h, ref, inject } from "vue";
import type { VNode } from "vue";
import type { ColumnType } from "ant-design-vue/es/table";
import { Table } from "ant-design-vue";
import { colors, Colors } from "./Colors.ts";
import Color from "colorjs.io";

export type ExtendedColumn<T = Record<string, any>> = ColumnType<T> & {
  extended?: boolean;
};

interface ExtendedRow<T> {
  originalRecord: T;
  columnIndex: number;
}

interface Props<T = Record<string, any>> {
  dataSource: T[];
  rowKey: string | ((record: T) => string);
  columns: ExtendedColumn<T>[];
  color: string;
  i18nPrefix: string;
}

const props = defineProps<Props>();

const theme = computed(() => {
  const mainColor = new Color(colors.value[Colors.Card]);

  const bgColor = mainColor;
  bgColor.alpha = 0.5;

  const black = new Color("black");
  const white = new Color("white");

  const contrastBlack = bgColor.contrast(black, "WCAG21");
  const contrastWhite = bgColor.contrast(white, "WCAG21");
  const contrastColor = contrastBlack > contrastWhite ? black : white;

  const textColor = mainColor.mix(contrastColor, 0.75, { space: "srgb" });

  return {
    token: {
      colorBgContainer: bgColor.toString(),
      colorTextHeading: textColor.toString(),
      colorText: textColor.toString(),
      algorithm: true,
    },
    components: {
      Table: {
        borderRadiusSM: 0,
        borderRadiusMD: 0,
        borderRadiusLG: 0,
        algorithm: true,
      },
    },
  };
});

const mainColumns = computed<ExtendedColumn[]>(() =>
  props.columns.filter((col) => !col.extended),
);

const extendedColumns = computed<ExtendedColumn[]>(() =>
  props.columns.filter((col) => col.extended),
);

const extendedTableColumns = computed<ColumnType<ExtendedRow<any>>[]>(() => [
  {
    key: "key",
    width: "16ch",
    customRender: function ({ record }): VNode {
      const column = extendedColumns.value[record.columnIndex];
      return h("span", {
        "data-i18n": `${props.i18nPrefix}.${column.key}`,
      });
    },
  },
  {
    key: "value",
    customRender: ({ record }): VNode =>
      h(Table, {
        class: "table-extension-cell",
        columns: [extendedColumns.value[record.columnIndex]],
        dataSource: [record.originalRecord],
        pagination: false,
        showHeader: false,
        bordered: false,
        size: "small",
      }),
  },
]);

function extendedDataSource<T>(record: T): ExtendedRow<T>[] {
  return extendedColumns.value.map((_, index) => ({
    originalRecord: record,
    columnIndex: index,
  }));
}
</script>
