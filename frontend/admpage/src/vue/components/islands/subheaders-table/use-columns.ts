import { computed } from "vue";
import { type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type Header } from "@/services/api/generated/types.gen";
import {
  makeConfigLine,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderActions } from "./rendering";

export function useSubHeadersTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<Header>[] = [
      {
        key: "key",
        dataIndex: ["Key"],
        fixed: "left",
      },
      {
        key: "value",
        dataIndex: ["Value"],
        customRender: ({ text }) => makeConfigLine(text),
      },
      {
        key: "actions",
        customRender: ({ record }) => renderActions(record),
        fixed: "right",
      },
    ];
    return i18nateColumns<Header>(`${i18nPrefix}.columns`, columns);
  });
}
