import { computed } from "vue";
import { type ExtendedColumn, type Props } from "./table-types.ts";

export function useTableColumns<T>(columns: ExtendedColumn<T>[]) {
  const mainCols = computed(() => columns.filter((c) => !c.extended));
  const extCols = computed(() => columns.filter((c) => c.extended));
  return { mainCols, extCols };
}
