import { client } from "./generated/client.gen";
import { config } from "@/config/config";

export function setupClient() {
  client.setConfig({
    baseUrl: config.ApiPrefix,
  });
}
