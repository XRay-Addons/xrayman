import type { App } from "vue";
import Antd from "ant-design-vue";
import PrimeVue from "primevue/config";
import Aura from "@primeuix/themes/aura";

export default (app: App) => {
  app.use(Antd);
  app.use(PrimeVue, {
    theme: {
      preset: Aura,
    },
  });
};
