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
      <!-- Expanded row -->
      <template #expandedRowRender="{ record }">
        <a-table
          :columns="extendedTableColumns"
          :data-source="extendedDataSource(record)"
          :pagination="false"
          :show-header="false"
          :showExpandColumn="false"
          size="small"
          class="table-extension"
        />
      </template>
    </a-table>
  </a-config-provider>
</template>

<script setup lang="ts" generic="T">
import { computed, h } from "vue";
import type { VNode } from "vue";
import type { ColumnType } from "ant-design-vue/es/table";
import { Table } from "ant-design-vue";
import { Colors, colors } from "../../state/colors";
import Color from "colorjs.io";

// template types
export type ExtendedColumn<T> = ColumnType<T> & {
  extended?: boolean;
};

interface ExtendedRow<T> {
  originalRecord: T;
  columnIndex: number;
}

interface Props<T> {
  dataSource: T[];
  rowKey: string | ((record: T) => string);
  columns: ExtendedColumn<T>[];
}

// props + template params
const props = defineProps<Props<T>>();

// theme
const theme = computed(() => {
  const mainColor = new Color(colors.value[Colors.Card]);

  const bgColor = mainColor.clone();
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

// columns
const mainColumns = computed<ExtendedColumn<T>[]>(() =>
  props.columns.filter((col) => !col.extended),
);

const extendedColumns = computed<ExtendedColumn<T>[]>(() =>
  props.columns.filter((col) => col.extended),
);

// expanded columns as rows
const extendedTableColumns = computed<ColumnType<ExtendedRow<T>>[]>(() => [
  {
    // show only header and hide the body
    key: "key",
    width: "16ch",
    customRender: ({ record }): VNode => {
      // unbold table header for better style
      const keyHeader = unbold(extendedColumns.value[record.columnIndex]);
      return h(Table, {
        class: "table-extension-key-cell no-body-table",
        style: { "margin-inline": "0" },
        columns: [keyHeader],
        dataSource: [],
        showHeader: true,
        bordered: false,
        pagination: false,
        showExpandColumn: false,
        indentSize: 0,
      });
    },
  },
  {
    // show only cell and hide the header
    key: "value",
    customRender: ({ record }): VNode => {
      return h(Table, {
        class: "table-extension-value-cell",
        style: { "margin-inline": "0" },
        columns: [extendedColumns.value[record.columnIndex]],
        dataSource: [record.originalRecord],
        pagination: false,
        showHeader: false,
        bordered: false,
        indentSize: 0,
        showExpandColumn: false,
        size: "small",
      });
    },
  },
]);

function extendedDataSource(record: T): ExtendedRow<T>[] {
  return extendedColumns.value.map((_, index) => ({
    originalRecord: record,
    columnIndex: index,
  }));
}

function unbold<T>(keyHeader: ExtendedColumn<T>): ExtendedColumn<T> {
  const originalFn = keyHeader.customHeaderCell;

  keyHeader.customHeaderCell = (...args) => {
    const originalProps = originalFn ? originalFn(...args) : {};
    return {
      ...originalProps,
      style: Object.assign({}, originalProps.style, {
        fontWeight: "normal",
        borderBottom: "none",
      }),
    };
  };

  return keyHeader;
}
</script>

<style>
.no-body-table tbody {
  display: none;
}
</style>
