import { t } from "@/runtime/i18n";
import { DynamicConfig } from "@/services/api/generated";
import { Platform, checkPlatform } from "@xrayman/shared/runtime/platforms/platforms";

export function checkConfig(cfg: DynamicConfig): Error | undefined {
  const platforms: Platform[] = ["ios", "macos", "android", "windows", "unknown"];

  for (const platform of platforms) {
    const count = cfg.AppLinks.filter((app) => checkPlatform(app.Platforms, platform)).length;

    if (count === 0) {
      const error = new Error(`${platform}: ${t("cards.config.errors.no-platform-apps")}`);
      error.name = t("cards.config.errors.app-links");
      return error;
    }
  }

  return undefined;
}
