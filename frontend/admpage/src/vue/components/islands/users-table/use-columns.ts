import { computed } from "vue";
import { sm, type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type UserView } from "@/services/api/generated/types.gen";
import {
  makeConfigLine,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import {
  renderTag,
  renderTraffic,
  renderActions,
  renderUserPageURL,
  renderApiUrl,
} from "./rendering";

export function useUsersTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<UserView>[] = [
      {
        key: "display-name",
        dataIndex: ["User", "Profile", "DisplayName"],
      },
      {
        key: "traffic-total",
        dataIndex: ["Traffic", "Total"],
        customRender: ({ value }) => renderTraffic(value),
        extended: sm,
      },
      {
        key: "traffic-last-month",
        dataIndex: ["Traffic", "LastMonth"],
        customRender: ({ value }) => renderTraffic(value),
        extended: sm,
      },
      {
        key: "target-status",
        dataIndex: ["User", "TargetStatus"],
        customRender: ({ value }) => renderTag(value),
      },
      {
        key: "name",
        dataIndex: ["User", "Profile", "Name"],
        extended: true,
      },
      {
        key: "id",
        dataIndex: ["User", "Profile", "ID"],
        extended: true,
      },
      {
        key: "vless-uuid",
        dataIndex: ["User", "Profile", "VlessUUID"],
        customRender: ({ text }) => makeConfigLine(text),
        extended: true,
      },
      {
        key: "userpage",
        dataIndex: ["User", "Profile"],
        customRender: ({ text }) => renderUserPageURL(text),
        extended: true,
      },
      {
        key: "subscription",
        dataIndex: ["User", "Profile", "SubscriptionPath"],
        customRender: ({ text }) => renderApiUrl(text),
        extended: true,
      },
      {
        key: "actions",
        dataIndex: ["User", "TargetStatus"],
        customRender: ({ value, record }) => renderActions(value, record.User),
        extended: true,
      },
    ];
    return i18nateColumns<UserView>(`${i18nPrefix}.columns`, columns);
  });
}
