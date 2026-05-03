import { client } from "./generated/client.gen";

export function setupClient() {
  client.setConfig({
    baseUrl: "http://localhost:80/api",
  });
}
