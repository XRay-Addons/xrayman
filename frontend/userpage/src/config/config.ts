import { notifyApiError } from "@/runtime/notifications/use-notifications";
import type { UserPageConfig } from "./config.d";

type CfgWindow = Window & {
  __CONFIG__: UserPageConfig;
};

export function config(): UserPageConfig {
  const cfg = (window as unknown as CfgWindow).__CONFIG__;
  if (!cfg) {
    notifyApiError("config_js");
    throw new Error("can't load config");
  }
  return cfg;
}
