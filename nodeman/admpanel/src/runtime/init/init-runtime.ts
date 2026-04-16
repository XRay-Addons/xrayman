import { initI18nObserver } from "./i18n";
import { fixIOSActive } from "./ios-active-fix";
import { serverPoll } from "@/runtime/transport/server-poll";

function initRuntime() {
  initI18nObserver();
  fixIOSActive();
  serverPoll.start();
}

initRuntime();
