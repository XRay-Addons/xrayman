import { notifyApiError } from "@/runtime/notifications/use-notifications";
import type { AdminPageConfig } from "./config.d";

type CfgWindow = Window & {
  __CONFIG__: AdminPageConfig;
};

//export const config = (window as unknown as CfgWindow).__CONFIG__;

export function config(): AdminPageConfig {
  const cfg = (window as unknown as CfgWindow).__CONFIG__;
  if (!cfg) {
    notifyApiError("config_js");
    throw new Error("can't load config");
  }
  return cfg;
}
