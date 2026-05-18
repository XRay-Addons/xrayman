import type { AdminPageConfig } from "./config.d";
import { notifyError } from "@/runtime/notifications/use-notifications";
import { makeSingleton } from "@xrayman/shared/runtime/singletone/singletone";

export const config = makeSingleton<AdminPageConfig>(async () => {
  const r = await fetch("./config.json");

  if (!r.ok) {
    notifyError("errors.server.config-json");
    throw new Error("config load failed");
  }

  return r.json();
});
