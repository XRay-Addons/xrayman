<template>
  <ExtendedTable
    :data-source="nodes"
    :columns="nodesColumns"
    :row-key="rowKey"
    :loading="nodesLoading"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
import ExtendedTable, {
  type ExtendedColumn,
} from "@/components/ui/TableExt.vue";
import type {
  Node as APINode,
  NodeStatus as APINodeStatus,
} from "@/api/generated/types.gen";
import { onMounted, type VNode, computed } from "vue";
import {
  i18nateColumns,
  makeMonospace,
  enabledTag,
  disabledTag,
  unknownTag,
  enableBtn,
  disableBtn,
  ensureDeleteBtn,
  mergeActionBtns,
} from "@/lib/table-ext-elements";
import { nodes, nodesLoading, reloadNodes } from "@/state/nodes";
import { enableUser, disableUser } from "@/api/client";

onMounted(reloadNodes);

// row key
const rowKey = (record: APINode): string => String(record.ID);

// i18n prefix
const i18nPrefix = "table.nodes";

// columns
const nodesColumns = computed(() => {
  const columns: ExtendedColumn<APINode>[] = [
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
    },
    {
      key: "target-status",
      dataIndex: ["TargetStatus"],
      customRender: ({ value }) => renderTag(value),
    },
    {
      key: "access-key",
      dataIndex: ["Config", "ConnectionInfo", "AccessKey"],
      customRender: ({ text }) => makeMonospace(text),
      ellipsis: true,
      width: "8ch",
      extended: true,
    },
    {
      key: "client-config",
      dataIndex: ["Config", "ClientConfigTemplate"],
      customRender: ({ text }) => makeMonospace(text),
      extended: true,
    },
    {
      key: "actions",
      dataIndex: ["TargetStatus"],
      customRender: ({ value }) => renderActions(value),
      extended: true,
    },
  ];

  return i18nateColumns<APINode>(`${i18nPrefix}.columns`, columns);
});

// value rendering
function renderTag(status: APINodeStatus) {
  if (status === "stopped") {
    return disabledTag("table.nodes.status.stopped");
  } else if (status === "running") {
    return enabledTag("table.nodes.status.running");
  } else {
    return unknownTag("table.nodes.status.unknown");
  }
}

function startNodeFn(node: APINode): BtnAction {
  return async () => {
    const r = await enableUser(-1);
    if (r.ok) {
      reloadNodes();
    } else {
      console.log(r.reason);
    }
  };
}

function stopNodeFn(node: APINode): BtnAction {
  return async () => {
    const r = await disableUser(-1);
    if (r.ok) {
      reloadNodes();
    } else {
      console.log(r.reason);
    }
  };
}

function renderActions(status: APINodeStatus, node: APINode) {
  const actions: VNode[] = [];

  if (status !== "running") {
    actions.push(enableBtn("table.nodes.actions.start", startNodeFn(node)));
  }
  if (status !== "stopped") {
    actions.push(disableBtn("table.nodes.actions.stop", stopNodeFn(node)));
  }
  actions.push(ensureDeleteBtn("table.nodes.actions"));

  return mergeActionBtns(actions);
}
</script>
