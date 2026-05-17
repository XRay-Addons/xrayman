import type { UserPageConfig } from "./config.d";
import { notifyError } from "@/runtime/notifications/use-notifications";

export const config: UserPageConfig = await fetch("./config.json").then(async (r) => {
  if (!r.ok) {
    notifyError("errors.server.config-json");
    throw new Error("config load failed");
  }
  return (await r.json()) as UserPageConfig;
});
