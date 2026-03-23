import PLATFORM_BUTTONS from "./platform-apps.json";

type Platform = keyof typeof PLATFORM_BUTTONS;

function mapPlatform(p: string): Platform {
  switch (p.toLowerCase()) {
    case "macos":
      return "macos";
    case "ios":
      return "ios";
    case "android":
      return "android";
    case "windows":
      return "windows";
    default:
      return "unknown";
  }
}

function detectPlatform(): Platform {
  const uaData = (navigator as any).userAgentData;
  if (uaData?.platform) {
    return mapPlatform(uaData.platform);
  }

  const ua = navigator.userAgent.toLowerCase();

  if (/iphone|ipad|ipod/.test(ua)) return "ios";
  if (/macintosh/.test(ua)) return "macos";
  if (/android/.test(ua)) return "android";
  if (/windows/.test(ua)) return "windows";

  return "unknown";
}

function setupInstallButtons() {
  const templateId = "install-app-btn-template";
  const template = document.getElementById(
    templateId,
  ) as HTMLTemplateElement | null;

  if (!template) {
    console.error(`Template "${templateId}" not found`);
    return;
  }

  const parent = template.parentElement;
  if (!parent) {
    console.error(`Template "${templateId}" has no parent`);
    return;
  }

  const platform = detectPlatform();
  const buttons = PLATFORM_BUTTONS[platform] ?? [];

  const frag = document.createDocumentFragment();
  buttons.forEach((data) => {
    const fragment = template.content.cloneNode(true) as DocumentFragment;
    const buttonEl = fragment.querySelector<HTMLButtonElement>("button");

    if (buttonEl) {
      buttonEl.addEventListener("click", () => {
        window.open(data.url, "_blank", "noopener");
      });
    }
    const textEl = fragment.querySelector<HTMLSpanElement>(
      ".install-app-btn-text",
    );
    if (textEl) {
      textEl.setAttribute("data-i18n", data.text);
    }

    frag.appendChild(fragment);
  });
  parent.appendChild(frag);
}

setupInstallButtons();
