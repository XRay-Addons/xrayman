import type { AdminPageConfig } from "./config.d";
import { notifyError } from "@/runtime/notifications/use-notifications";

export const config: AdminPageConfig = await fetch("./config.json").then(async (r) => {
  if (!r.ok) {
    notifyError("errors.server.config-json");
    throw new Error("config load failed");
  }
  return (await r.json()) as AdminPageConfig;
});
