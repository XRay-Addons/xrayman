import type { ColumnType } from "ant-design-vue/es/table";

export type ExtendedColumn<T> = ColumnType<T> & {
  extended?: boolean;
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
