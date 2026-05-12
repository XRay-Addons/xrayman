import type { ExtendedColumn, ExtendedRow } from "./table-types";
import { ConfigProvider, Table } from "ant-design-vue";
import { h, computed, type VNode } from "vue";

type ExpandRenderCtx<T> = {
  record: T;
  index: number;
  indent: number;
  expanded: boolean;
};

export function useExpandedRowRender<T>(
  extendedColumns: ExtendedColumn<T>[],
  keyColumnWidth: string,
) {
  if (extendedColumns.length == 0) {
    return null;
  }

  function extendedDataSource(record: T): ExtendedRow<T>[] {
    const r = extendedColumns.map((_, index) => ({
      originalRecord: record,
      columnIndex: index,
    }));
    return r;
  }

  // ellipsis required to make cells scrollable (wtf framework)
  const extendedTableColumns = [
    {
      // show only header and hide the body
      key: "key",
      width: keyColumnWidth,
      customRender: keyCellRender(extendedColumns),
      ellipsis: true,
    },
    {
      // show only cell and hide the header
      key: "value",
      customRender: valCellRender(extendedColumns),
      ellipsis: true,
    },
  ];

  return (ctx: ExpandRenderCtx<T>) => {
    return h(
      ConfigProvider,
      {
        theme: expandedTableTheme,
      },

      () =>
        h(Table, {
          dataSource: extendedDataSource(ctx.record),
          columns: extendedTableColumns,
          class: "table-ext-expand-table",
          scroll: { x: "100%" },
          showHeader: false,
          bordered: false,
          pagination: false,
          showExpandColumn: false,
          indentSize: 0,
        }),
    );
  };
}

type CellRenderCtx<T> = {
  text: any;
  record: ExtendedRow<T>;
  index: number;
  column: ExtendedColumn<T>;
  value: any;
};

const expandedTableTheme = {
  components: {
    Table: {
      colorFillAlter: "transparent",
      algorithm: false,
      fontWeightStrong: 400, // normal
    },
  },
};

function keyCellRender<T>(columns: ExtendedColumn<T>[]) {
  return (ctx: CellRenderCtx<T>): VNode => {
    const keyHeader = columns[ctx.record.columnIndex];
    return h(Table, {
      columns: [keyHeader],
      dataSource: [],

      class: "table-ext-key-cell",
      style: { "margin-inline": "0" },
      size: "small",
      showHeader: true,
      bordered: false,
      pagination: false,
      showExpandColumn: false,
      indentSize: 0,
    });
  };
}

function valCellRender<T>(columns: ExtendedColumn<T>[]) {
  return (ctx: CellRenderCtx<T>): VNode => {
    return h(Table, {
      columns: [columns[ctx.record.columnIndex]],
      dataSource: [ctx.record.originalRecord],

      class: "table-ext-value-cell",
      style: { "margin-inline": "0" },
      scroll: { x: "max-content" },
      size: "small",
      pagination: false,
      showHeader: false,
      bordered: false,
      indentSize: 0,
      showExpandColumn: false,
    });
  };
}

/*
          dataSource: extendedDataSource(ctx.record),
          columns: extendedTableColumns,

          class: "table-ext-expand-table",
          showHeader: false,
          bordered: false,
          pagination: false,
          showExpandColumn: false,
          indentSize: 0,
        }),
*/
