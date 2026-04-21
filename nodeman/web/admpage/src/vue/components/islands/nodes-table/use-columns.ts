import { computed } from "vue";
import { type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type Node } from "@/services/api/generated/types.gen";
import {
  makeCopyable,
  makeMonospace,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderTag, renderActions } from "./rendering";

export function useNodesTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<Node>[] = [
      {
        key: "id",
        dataIndex: ["ID"],
        width: "8ch",
      },
      {
        key: "endpoint",
        dataIndex: ["Config", "ConnectionInfo", "Endpoint"],
        width: "24ch",
      },
      {
        key: "current-status",
        dataIndex: ["CurrentStatus"],
        customRender: ({ value }) => renderTag(value),
        width: "16ch",
      },
      {
        key: "target-status",
        dataIndex: ["TargetStatus"],
        customRender: ({ value }) => renderTag(value),
        width: "16ch",
      },
      {
        key: "access-key",
        dataIndex: ["Config", "ConnectionInfo", "AccessKey"],
        customRender: ({ text }) => makeCopyable(makeMonospace(text), text),
        ellipsis: true,
        width: "8ch",
        extended: true,
      },
      {
        key: "client-config",
        dataIndex: ["Config", "ClientConfigTemplate"],
        customRender: ({ text }) => {
          return makeMonospace(JSON.stringify(text));
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
