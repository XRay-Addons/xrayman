import { MakeFullUrl } from "@xrayman/shared/runtime/paths/paths";
import { config } from "@/config/config";

export function MakeApiUrl(path: string): string {
  return MakeFullUrl(config.ApiPrefix, path);
}

export function MakeUserpageURL(id: number, name: string): string {
  return MakeFullUrl(config.UserPrefix, `${id}-${name}`);
}
