import { h, type VNode } from "vue";
import { Tag, Button, Popconfirm, Space } from "ant-design-vue";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons-vue";

export function enabledTag(i18n: string) {
  return makeTag("success", i18n, CheckCircleOutlined);
}

export function disabledTag(i18n: string) {
  return makeTag("error", i18n, CloseCircleOutlined);
}

export function unknownTag(i18n: string) {
  return makeTag("warning", i18n, ExclamationCircleOutlined);
}

function makeTag(color: string, i18n: string, icon: any): VNode {
  return h(
    Tag,
    { color },
    {
      default: () => [h(icon), h("span", { "data-i18n": i18n })],
    },
  );
}

export function enableBtn(i18n: string): VNode {
  return h(Button, {
    ghost: true,
    size: "small",
    type: "primary",
    "data-i18n": i18n,
  });
}

export function disableBtn(i18n: string): VNode {
  return h(Button, {
    danger: true,
    ghost: true,
    size: "small",
    type: "primary",
    "data-i18n": i18n,
  });
}

export function ensureDeleteBtn(i18nPrefix: string): VNode {
  return h(
    Popconfirm,
    {
      okText: h("span", {
        "data-i18n": `${i18nPrefix}.delete-confirn-yes`,
      }),
      cancelText: h("span", {
        "data-i18n": `${i18nPrefix}.delete-confirn-no`,
      }),
    },
    {
      default: () =>
        h(Button, {
          danger: true,
          size: "small",
          type: "primary",
          style: { boxShadow: "none" },
          "data-i18n": `${i18nPrefix}.delete`,
        }),
      title: () =>
        h("span", {
          "data-i18n": `${i18nPrefix}.delete-confirn-header`,
        }),
      description: () =>
        h("span", {
          "data-i18n": `${i18nPrefix}.delete-confirn-body`,
        }),
    },
  );
}

export function mergeActionBtns(btns: VNode[]): VNode {
  return h(Space, { size: "small" }, () => btns);
}
