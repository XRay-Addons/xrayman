import { computed } from "vue";
import { type ExtendedColumn } from "./table-types";

export function useTableColumns<T>(columns: ExtendedColumn<T>[]) {
  const isExtended = (column: ExtendedColumn<T>) => {
    const value = column.extended;
    if (typeof value === "function") {
      return value();
    }
    return value === true;
  };

  const mainCols = computed(() => columns.filter((c) => !isExtended(c)));
  const extCols = computed(() => columns.filter(isExtended));

  return { mainCols, extCols };
}
