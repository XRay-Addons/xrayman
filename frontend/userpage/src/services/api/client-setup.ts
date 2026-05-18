import { client } from "./generated/client.gen";
import { config } from "@/config/config";
import { makeSingleton } from "@xrayman/shared/runtime/singletone/singletone";

export const clientSetup = makeSingleton<void>(async () => {
  client.setConfig({
    baseUrl: (await config.get()).routes.api_prefix,
  });
});
