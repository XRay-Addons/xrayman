import { Notyfier } from "@xrayman/shared/runtime/notifications/notyfier";
import type { ApiReason } from "@xrayman/shared/services/api/api-reason";
import { t } from "@/runtime/i18n";

const ntf = new Notyfier();

export function notifyError(message: string, details?: string) {
  ntf.errorNotification(t(message), details && t(details));
}

export function notifyApiError(errkey: string, reason?: ApiReason) {
  ntf.errorNotification(
    t(`errors.server.${errkey}`),
    reason && t(`errors.server.reason.${reason}`),
  );
}
