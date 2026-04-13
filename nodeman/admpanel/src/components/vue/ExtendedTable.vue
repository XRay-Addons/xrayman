<template>
  <a-config-provider :theme="theme">
    <a-table
      :data-source="dataSource"
      :columns="mainColumns"
      :row-key="rowKey"
      :scroll="{ x: 'max-width' }"
      size="medium"
      class="table-extended"
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
//import { Colors } from "./Globals.ts";

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

/*const themeData = computed(() => ({
  token: {
    colorBgContainer: "red",
    algorythm: true,
    algorithm: true,
  },
  components: {
    Table: {
      borderRadiusSM: 0,
      borderRadiusMD: 0,
      borderRadiusLG: 0,
      algorithm: true,
    },
    Button: {
      borderRadiusSM: 4,
      borderRadiusMD: 4,
      borderRadiusLG: 4,
    },
  },
}));

const theme = ref(themeData);*/

const theme = computed(() => ({
  token: {
    colorBgContainer: colors.value[Colors.Card],
    algorythm: true,
    algorithm: true,
    colorText: "white",
  },
  components: {
    Table: {
      borderRadiusSM: 0,
      borderRadiusMD: 0,
      borderRadiusLG: 0,
      algorithm: true,
      headerColor: "white",
    },
    Button: {
      borderRadiusSM: 4,
      borderRadiusMD: 4,
      borderRadiusLG: 4,
    },
  },
}));

/*setTimeout(() => {
  theme.value.token.colorBgContainer = "#ff000080";
}, 2000);

setTimeout(() => {
  theme.value.token.colorBgContainer = "#33ff0080";
  console.log(theme.value);
}, 2000);*/

/*const theme = computed(() => ({
token: {
        colorBgContainer: props.color,
      },
      components: {
        Table: {
          borderRadiusSM: 0,
          borderRadiusMD: 0,
          borderRadiusLG: 0,
          algorithm: true,
        },
        Button: {
          borderRadiusSM: 4,
          borderRadiusMD: 4,
          borderRadiusLG: 4,
        },
      },
    }"
}));*/
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
