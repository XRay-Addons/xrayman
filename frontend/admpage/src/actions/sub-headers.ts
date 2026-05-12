import { listSubHeaders } from "@/services/api/client";
import { subHeaders, subHeadersLoading } from "@/state/sub-headers";
import { notifyApiError } from "@/runtime/notifications/use-notifications";

export async function reloadSubHeaders({
  quiet = false,
}: {
  quiet?: boolean;
} = {}) {
  if (quiet) {
    const result = await listSubHeaders();
    if (result.ok) subHeaders.value = result.data;
    return;
  }

  subHeadersLoading.value = true;

  try {
    const result = await listSubHeaders();

    if (result.ok) {
      subHeaders.value = result.data;
    } else {
      notifyApiError("get_sub_headers", result.reason);
    }
  } finally {
    subHeadersLoading.value = false;
  }
}
