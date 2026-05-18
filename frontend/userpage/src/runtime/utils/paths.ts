import { MakeFullUrl } from "@xrayman/shared/runtime/paths/paths";
import { config } from "@/config/config";

export async function MakeApiUrl(path: string): Promise<string> {
  return MakeFullUrl((await config.get()).routes.api_prefix, path);
}

export async function MakePageUrl(path: string): Promise<string> {
  return MakeFullUrl((await config.get()).routes.user_prefix, path);
}
