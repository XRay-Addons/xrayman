import { getNodeManagementAPI } from "./generated";
import { UserID, UserAPIData, ApiReason, API_REASON } from "../lib/types";
import { config } from "../config/config";
import axios from "axios";

export const api = axios.create({
  baseURL: config.API_URLPATH,
  validateStatus: () => true,
});

export const nodeApi = getNodeManagementAPI(api);

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; reason: ApiReason };

export async function newUser(
  DisplayName: string,
): Promise<ApiResult<UserAPIData>> {
  return handleApi(
    () => nodeApi.newUser({ DisplayName }),
    (data) => ({
      id: data.Profile.ID,
      name: data.Profile.Name,
      displayName: data.Profile.DisplayName,
    }),
  );
}

export async function getUser(userID: UserID): Promise<ApiResult<UserAPIData>> {
  return handleApi(
    () => nodeApi.getUser(userID.id, userID.name),
    (data) => ({
      id: data.Profile.ID,
      name: data.Profile.Name,
      displayName: data.Profile.DisplayName,
    }),
  );
}

async function handleApi<TData, TResult>(
  request: () => Promise<{ status: number; data?: TData }>,
  onSuccess: (data: TData) => TResult,
): Promise<ApiResult<TResult>> {
  try {
    const res = await request();

    if (res.status === 200 && res.data) {
      return { ok: true, data: onSuccess(res.data) };
    }

    switch (res.status) {
      case 400:
        return { ok: false, reason: API_REASON.BAD_REQUEST };
      case 401:
        return { ok: false, reason: API_REASON.UNAUTHORIZED };
      case 403:
        return { ok: false, reason: API_REASON.FORBIDDEN };
      case 404:
        return { ok: false, reason: API_REASON.NOT_FOUND };
      default:
        return { ok: false, reason: API_REASON.UNKNOWN };
    }
  } catch {
    return { ok: false, reason: API_REASON.NETWORK };
  }
}
