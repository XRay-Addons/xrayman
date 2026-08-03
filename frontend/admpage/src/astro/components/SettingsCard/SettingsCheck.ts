import { t } from "@/runtime/i18n";
import { Settings } from "@/services/api/generated";
import { Platform, checkPlatform } from "@xrayman/shared/runtime/platforms/platforms";

export function checkSettings(settings: Settings): Error | undefined {
  const platforms: Platform[] = ["ios", "macos", "android", "windows", "unknown"];

  for (const platform of platforms) {
    const count = settings.AppLinks.filter((app) => checkPlatform(app.Platforms, platform)).length;

    if (count === 0) {
      const error = new Error(`${platform}: ${t("cards.settings.errors.no-platform-apps")}`);
      error.name = t("cards.settings.errors.app-links");
      return error;
    }
  }

  return undefined;
}
