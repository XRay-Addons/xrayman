import { t } from "@/runtime/i18n";
import { type ApiReason } from "@xrayman/shared/services/api/api-reason";

type NotificationFn = (message: string, description?: string) => void;

let errorHandler: NotificationFn | null = null;

export function setErrorHandler(handler: NotificationFn | null) {
  errorHandler = handler;
}

function errorNotification(message: string, description?: string) {
  if (errorHandler) {
    errorHandler(message, description);
  } else {
    console.error(message, description);
  }
}

export function notifyApiError(errkey: string, reason?: ApiReason) {
  console.log(errkey);
  errorNotification(
    t(`errors.server.${errkey}`),
    reason ? t(`errors.server.reason.${reason}`) : "",
  );
}
