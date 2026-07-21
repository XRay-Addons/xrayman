import { client } from "./generated/client.gen";
import { config } from "@/config/config";

let initialized = false;

export function ensureClient() {
  if (initialized) return;

  client.setConfig({
    baseUrl: config().routes.api_prefix,
  });

  initialized = true;
}
