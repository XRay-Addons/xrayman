import type { App } from "vue";
import Antd from "ant-design-vue";
import { i18n } from "@/runtime/i18n/i18n";

export default (app: App) => {
  app.use(i18n);
  app.use(Antd);
};
