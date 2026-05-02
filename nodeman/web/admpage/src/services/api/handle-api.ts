import { API_REASON, type ApiReason } from "./api-reason";
import type { Error } from "./generated/types.gen";
import { authMan } from "./auth-man";

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
      return { ok: true, data: transform(resp.data) };
    }
    console.log("api call error:", resp.error);
    return { ok: false, reason: getErrorReason(resp.response.status) };
  } catch (error) {
    console.log("api call error:", error);
    return { ok: false, reason: API_REASON.NETWORK };
  }
}

function getErrorReason(status: number): ApiReason {
  switch (status) {
    case 400:
      return API_REASON.BAD_REQUEST;
    case 401:
      return API_REASON.UNAUTHORIZED;
    case 403:
      return API_REASON.FORBIDDEN;
    case 404:
      return API_REASON.NOT_FOUND;
    default:
      return API_REASON.UNKNOWN;
  }
}
