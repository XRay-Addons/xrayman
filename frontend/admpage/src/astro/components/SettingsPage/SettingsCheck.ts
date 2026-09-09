import { t } from "@/runtime/i18n";
import { Settings } from "@/services/api/generated";
import { Platform, checkPlatform } from "@xrayman/shared/runtime/platforms/platforms";

export function checkSettings(settings: Settings): Error | undefined {
  if (settings.RecentDays <= 0) {
    const error = new Error(`${t("settings.errors.positive-recent-days")}`);
    error.name = t("settings.errors.recent-days");
    return error;
  }
  if (settings.UpdateInterval <= 0) {
    const error = new Error(`${t("settings.errors.positive-update-interval")}`);
    error.name = t("settings.errors.update-interval");
    return error;
  }

  const platforms: Platform[] = ["ios", "macos", "android", "windows", "unknown"];
  for (const platform of platforms) {
    const count = settings.AppLinks.filter((app) => checkPlatform(app.Platforms, platform)).length;

    if (count === 0) {
      const error = new Error(`${platform}: ${t("settings.errors.no-platform-apps")}`);
      error.name = t("settings.errors.app-links");
      return error;
    }
  }

  return undefined;
}
