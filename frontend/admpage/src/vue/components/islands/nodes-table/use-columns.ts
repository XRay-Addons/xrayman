import { computed } from "vue";
import { sm, type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type Node } from "@/services/api/generated/types.gen";
import {
  makeConfigLine,
  makeConfigText,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderTag, renderActions } from "./rendering";

export function useNodesTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<Node>[] = [
      {
        key: "endpoint",
        dataIndex: ["Config", "ConnectionInfo", "Endpoint"],
      },
      {
        key: "current-status",
        dataIndex: ["CurrentStatus"],
        customRender: ({ value }) => renderTag(value),
        extended: sm,
      },
      {
        key: "target-status",
        dataIndex: ["TargetStatus"],
        customRender: ({ value }) => renderTag(value),
        extended: sm,
      },
      {
        key: "id",
        dataIndex: ["ID"],
        width: "10%",
        extended: true,
      },
      {
        key: "version",
        dataIndex: ["Config", "Settings", "Version"],
        customRender: ({ text }) => {
          return makeConfigLine(text);
        },
        extended: true,
      },
      {
        key: "access-key",
        dataIndex: ["Config", "ConnectionInfo", "AccessKey"],
        customRender: ({ text }) => makeConfigLine(text),
        width: "8ch",
        extended: true,
      },
      {
        key: "client-config",
        dataIndex: ["Config", "Settings", "ClientConfigTemplate"],
        customRender: ({ text }) => {
          return makeConfigText(JSON.stringify(text, null, 2));
        },
        extended: true,
      },
      {
        key: "actions",
        dataIndex: ["TargetStatus"],
        customRender: ({ value, record }) => renderActions(value, record),
        extended: true,
      },
    ];
    return i18nateColumns<Node>(`${i18nPrefix}.columns`, columns);
  });
}
