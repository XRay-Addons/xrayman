import { t } from "@/runtime/i18n/i18n";
import { API_REASON } from "@/lib/types";

type NotificationFn = (message: string, description: string) => void;

let errorHandler: NotificationFn | null = null;

export function setErrorHandler(handler: NotificationFn) {
  errorHandler = handler;
}

export function errorNotification(message: string, description: string) {
  if (errorHandler) {
    errorHandler(message, description);
  } else {
    console.error(message, description);
  }
}

export function serverErrorNotification(errkey: string, reason: API_REASON) {
  console.log(errkey);
  errorNotification(
    t(`errors.server.${errkey}`),
    t(`errors.server.reason.${reason}`),
  );
}
