import { computed } from "vue";
import { type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type Header } from "@/services/api/generated/types.gen";
import {
  makeCopyable,
  makeMonospace,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderActions } from "./rendering";

export function useSubHeadersTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<Header>[] = [
      {
        key: "id",
        dataIndex: ["ID"],
        width: "8ch",
      },
      {
        key: "key",
        dataIndex: ["Key"],
        ellipsis: true,
        width: "16ch",
      },
      {
        key: "value",
        dataIndex: ["Value"],
      },
      {
        key: "actions",
        customRender: ({ record }) => renderActions(record),
      },
    ];
    return i18nateColumns<Header>(`${i18nPrefix}.columns`, columns);
  });
}
