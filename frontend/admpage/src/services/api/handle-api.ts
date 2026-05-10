import {
  API_REASON,
  getErrorReason,
  type ApiReason,
} from "@xrayman/shared/services/api/api-reason";
import type { Error } from "./generated/types.gen";

export type ApiResult<T> = { ok: true; data: T } | { ok: false; reason: ApiReason };

export type ApiResponse<T> = ({ data: T; error: undefined } | { data: undefined; error: Error }) & {
  request: Request;
  response: Response;
};

export async function handleAPI<T, R>(
  apiCall: () => Promise<ApiResponse<T>>,
  transform: (data: T) => R,
): Promise<ApiResult<R>> {
  try {
    console.log("call api");
    let resp = await apiCall();
    if (!resp.error) {
      return { ok: true, data: transform(resp.data!) };
    }
    console.log("api call error:", resp.error);
    return { ok: false, reason: getErrorReason(resp.response.status) };
  } catch (error) {
    console.log("api call error:", error);
    return { ok: false, reason: API_REASON.NETWORK };
  }
}
