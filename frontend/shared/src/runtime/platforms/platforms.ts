export type Platform = "ios" | "macos" | "android" | "windows" | "unknown";

export function detectPlatform(): Platform {
  const ua = navigator.userAgent.toLowerCase();
  if (/iphone|ipad|ipod/.test(ua)) return "ios";
  if (/macintosh/.test(ua)) return "macos";
  if (/android/.test(ua)) return "android";
  if (/windows/.test(ua)) return "windows";
  return "unknown";
}

export function checkPlatform(regex: string, platform: Platform): boolean {
  try {
    return new RegExp(regex, "i").test(platform);
  } catch {
    return false;
  }
}
