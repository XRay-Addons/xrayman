import type { AdminPageConfig } from "./config.d";

type CfgWindow = Window & {
  __CONFIG__: AdminPageConfig;
};

export const config = (window as unknown as CfgWindow).__CONFIG__;
