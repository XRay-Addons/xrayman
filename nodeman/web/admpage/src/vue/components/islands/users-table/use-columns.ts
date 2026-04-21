import { computed } from "vue";
import { type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type User } from "@/services/api/generated/types.gen";
import {
  makeCopyable,
  makeMonospace,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderTag, renderActions } from "./rendering";

export function useUsersTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<User>[] = [
      {
        key: "id",
        dataIndex: ["Profile", "ID"],
        width: "4ch",
      },
      {
        key: "display-name",
        dataIndex: ["Profile", "DisplayName"],
        ellipsis: true,
        width: "16ch",
      },
      {
        key: "target-status",
        dataIndex: ["TargetStatus"],
        customRender: ({ value }) => renderTag(value),
        width: "8ch",
      },
      {
        key: "name",
        dataIndex: ["Profile", "Name"],
        extended: true,
      },
      {
        key: "vless-uuid",
        dataIndex: ["Profile", "VlessUUID"],
        ellipsis: true,
        width: "8ch",
        customRender: ({ text }) => makeCopyable(makeMonospace(text), text),
        extended: true,
      },
      {
        key: "actions",
        dataIndex: ["TargetStatus"],
        customRender: ({ value, record }) => renderActions(value, record),
        extended: true,
      },
    ];
    return i18nateColumns<User>(`${i18nPrefix}.columns`, columns);
  });
}
