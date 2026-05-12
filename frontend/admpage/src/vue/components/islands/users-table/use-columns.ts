import { computed } from "vue";
import { type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type User } from "@/services/api/generated/types.gen";
import {
  makeConfigLine,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderTag, renderActions, renderUserPageURL, renderApiUrl } from "./rendering";
import { MakeUserpageURL } from "@/runtime/utils/paths";

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
        extended: true,
      },
      {
        key: "vless-uuid",
        dataIndex: ["Profile", "VlessUUID"],
        customRender: ({ text }) => makeConfigLine(text),
        extended: true,
      },
      {
        key: "userpage",
        dataIndex: ["Profile"],
        customRender: ({ text }) => renderUserPageURL(text),
        extended: true,
      },
      {
        key: "subscription",
        dataIndex: ["Profile", "SubscriptionPath"],
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
