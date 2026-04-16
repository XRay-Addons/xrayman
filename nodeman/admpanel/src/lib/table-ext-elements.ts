import { h, type VNode } from "vue";
import {
  Tag,
  Button,
  Popconfirm,
  Space,
  TypographyText,
  Typography,
} from "ant-design-vue";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons-vue";
import { type ExtendedColumn } from "@/components/ui/TableExt.vue";
import { i18n } from "@/runtime/i18n/i18n";

function t(text: string) {
  return {
    default: () => {
      return i18n.global.t(text);
    },
  };
}

export function i18nateColumns<T>(
  i18nPrefix: string,
  columns: ExtendedColumn<T>[],
): ExtendedColumn<T>[] {
  return columns.map((col) => ({
    ...col,
    title: col.title || i18n.global.t(`${i18nPrefix}.${col.key}`),
  }));
}

export function makeMonospace(text: string): VNode {
  return h(
    TypographyText,
    {
      style: { fontFamily: "monospace" },
    },
    () => text,
  );
}

export function makeCopyable(node: VNode, textToCopy: string): VNode {
  return h(
    TypographyText,
    {
      copyable: {
        text: textToCopy,
        tooltip: false,
      },
    },
    () => {
      return node;
    },
  );
}

export function enabledTag(i18n: string): VNode {
  return makeTag("success", i18n, CheckCircleOutlined);
}

export function disabledTag(i18n: string): VNode {
  return makeTag("error", i18n, CloseCircleOutlined);
}

export function unknownTag(i18n: string): VNode {
  return makeTag("warning", i18n, ExclamationCircleOutlined);
}

function makeTag(color: string, i18n: string, icon: any): VNode {
  return h(Tag, { color }, t(i18n));
}
export type BtnAction = () => void | Promise<void>;

export function enableBtn(i18n: string, onClick?: BtnAction): VNode {
  return h(
    Button,
    {
      ghost: true,
      size: "small",
      type: "primary",
      onClick: onClick,
    },
    t(i18n),
  );
}

export function disableBtn(i18n: string, onClick?: BtnAction): VNode {
  return h(
    Button,
    {
      danger: true,
      ghost: true,
      size: "small",
      type: "primary",
      onClick: onClick,
    },
    t(i18n),
  );
}

export function ensureDeleteBtn(i18nPrefix: string): VNode {
  return h(
    Popconfirm,
    {
      okText: h("span", {}, t(`${i18nPrefix}.delete.confirm.ok`)),
      cancelText: h("span", {}, t(`${i18nPrefix}.delete.confirm.cancel`)),
    },
    {
      default: () =>
        h(
          Button,
          {
            danger: true,
            size: "small",
            type: "primary",
            style: { boxShadow: "none" },
          },
          t(`${i18nPrefix}.delete.button`),
        ),
      title: () => h("span", {}, t(`${i18nPrefix}.delete.confirm.header`)),
      description: () => h("span", {}, t(`${i18nPrefix}.delete.confirm.body`)),
    },
  );
}

export function mergeActionBtns(btns: VNode[]): VNode {
  return h(Space, { size: "small" }, () => btns);
}
