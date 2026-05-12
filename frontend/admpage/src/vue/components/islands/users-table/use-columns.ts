import { computed } from "vue";
import { type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type User } from "@/services/api/generated/types.gen";
import {
  makeCopyable,
  makeMonospace,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderTag, renderActions, renderApiUrl } from "./rendering";

export function useUsersTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<User>[] = [
      {
        key: "display-name",
        dataIndex: ["Profile", "DisplayName"],
        width: "75%",
      },
      {
        key: "target-status",
        dataIndex: ["TargetStatus"],
        customRender: ({ value }) => renderTag(value),
        width: "25%",
      },
      {
        key: "name",
        dataIndex: ["Profile", "Name"],
        extended: true,
      },
      {
        key: "id",
        dataIndex: ["Profile", "ID"],
        width: "8ch",
        extended: true,
      },
      {
        key: "vless-uuid",
        dataIndex: ["Profile", "VlessUUID"],
        width: "8ch",
        customRender: ({ text }) => makeCopyable(makeMonospace(text), text),
        extended: true,
      },
      {
        key: "subscription",
        dataIndex: ["Profile", "SubscriptionPath"],
        width: "16ch",
        customRender: ({ text }) => renderApiUrl(text),
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
