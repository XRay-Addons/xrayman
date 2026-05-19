import type { UserPageConfig } from "./config.d";

type CfgWindow = Window & {
  __CONFIG__: UserPageConfig;
};

export const config = (window as unknown as CfgWindow).__CONFIG__;
