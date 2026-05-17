import { MakeFullUrl } from "@xrayman/shared/runtime/paths/paths";
import { config } from "@/config/config";

export function MakeApiUrl(path: string): string {
  return MakeFullUrl(config.routes.api_prefix, path);
}

export function MakeUserpageURL(id: number, name: string): string {
  return MakeFullUrl(config.routes.user_prefix, `${id}-${name}`);
}
