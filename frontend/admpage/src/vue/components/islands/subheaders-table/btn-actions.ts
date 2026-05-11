import { deleteSubHeader } from "@/services/api/client";
import { reloadSubHeaders } from "@/actions/sub-headers";
import { type Header } from "@/services/api/generated/types.gen";
import { type BtnAction } from "../../primitives/table-ext/render-primitives";
import { notifyApiError } from "@/runtime/notifications/use-notifications";

export function deleteHeaderAction(header: Header): BtnAction {
  return async () => {
    const r = await deleteSubHeader(header.ID);
    if (r.ok) {
      reloadSubHeaders({ quiet: true });
    } else {
      notifyApiError("delete_sub_header", r.reason);
    }
  };
}
