import { computed } from "vue";
import { sm, type ExtendedColumn } from "@/vue/components/primitives/table-ext/table-types";
import { type NodeView } from "@/services/api/generated/types.gen";
import {
  makeConfigLine,
  makeConfigText,
  i18nateColumns,
} from "@/vue/components/primitives/table-ext/render-primitives";
import { renderTag, renderActions, renderTraffic, renderLoad } from "./rendering";

export function useNodesTableColumns(i18nPrefix: string) {
  return computed(() => {
    const columns: ExtendedColumn<NodeView>[] = [
      {
        key: "endpoint",
        dataIndex: ["Node", "Config", "ConnectionInfo", "Endpoint"],
      },
      {
        key: "traffic-total",
        dataIndex: ["Traffic", "Total"],
        customRender: ({ value }) => renderTraffic(value),
        extended: sm,
      },
      {
        key: "traffic-recent-days",
        dataIndex: ["Traffic", "RecentDays"],
        customRender: ({ value }) => renderTraffic(value),
        extended: sm,
      },
      {
        key: "current-status",
        dataIndex: ["Node", "CurrentStatus"],
        customRender: ({ value }) => renderTag(value),
        extended: false,
      },
      {
        key: "target-status",
        dataIndex: ["Node", "TargetStatus"],
        customRender: ({ value }) => renderTag(value),
        extended: sm,
      },
      {
        key: "id",
        dataIndex: ["Node", "ID"],
        extended: true,
      },
      {
        key: "version",
        dataIndex: ["Node", "Config", "Settings", "Version"],
        customRender: ({ text }) => makeConfigLine(text),
        extended: true,
      },
      {
        key: "access-key",
        dataIndex: ["Node", "Config", "ConnectionInfo", "AccessKey"],
        customRender: ({ text }) => makeConfigLine(text),
        extended: true,
      },
      {
        key: "connections",
        dataIndex: ["Performance", "OpenConnections"],
        extended: true,
      },
      {
        key: "cpu-load",
        dataIndex: ["Performance", "CpuLoad"],
        customRender: ({ text }) => renderLoad(text),
        extended: true,
      },
      {
        key: "ram-load",
        dataIndex: ["Performance", "RamLoad"],
        customRender: ({ text }) => renderLoad(text),
        extended: true,
      },
      {
        key: "mem-load",
        dataIndex: ["Performance", "MemLoad"],
        customRender: ({ text }) => renderLoad(text),
        extended: true,
      },
      {
        key: "client-config",
        dataIndex: ["Node", "Config", "Settings", "ClientConfigTemplate"],
        customRender: ({ text }) => makeConfigText(JSON.stringify(text, null, 2)),
        extended: true,
      },
      {
        key: "actions",
        dataIndex: ["Node", "TargetStatus"],
        customRender: ({ value, record }) => renderActions(value, record.Node),
        extended: true,
      },
    ];
    return i18nateColumns<NodeView>(`${i18nPrefix}.columns`, columns);
  });
}
