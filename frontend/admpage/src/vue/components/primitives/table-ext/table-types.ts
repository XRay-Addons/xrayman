import type { ColumnType } from "ant-design-vue/es/table";

import { ref } from "vue";

const windowWidth = ref(typeof window !== "undefined" ? window.innerWidth : 0);

if (typeof window !== "undefined") {
  window.addEventListener("resize", () => {
    windowWidth.value = window.innerWidth;
  });
}

export function sm(): boolean {
  return windowWidth.value <= 576;
}

export function md(): boolean {
  return windowWidth.value <= 768;
}

export function xl(): boolean {
  return windowWidth.value <= 1200;
}

export function xxl(): boolean {
  return windowWidth.value <= 1400;
}

export type ExtendedColumn<T> = ColumnType<T> & {
  // flag to hide column in the on-demand expanded part of row.
  // extended === undefined || extended === false - always main
  // extended === true - always extended
  // extended === func - extended if func(width) === true
  extended?: boolean | (() => boolean);
};

export interface Props<T> {
  dataSource: T[];
  rowKey: string | ((record: T) => string);
  columns: ExtendedColumn<T>[];
  expandedKeyColumnWidth?: string;
}

export interface ExtendedRow<T> {
  originalRecord: T;
  columnIndex: number;
}
