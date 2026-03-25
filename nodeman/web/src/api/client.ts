import { getNodeManagementAPI } from "./generated";
import { UserID, UserAPIData } from "../lib/types";
import { config } from "../config/config";
import axios from "axios";

export const api = axios.create({
  baseURL: config.API_URLPATH,
});

export const nodeApi = getNodeManagementAPI(api);

export const API_REASON = {
  BAD_REQUEST: "bad_request",
  NOT_FOUND: "not_found",
  UNAUTHORIZED: "unauthorized",
  FORBIDDEN: "forbidden",
  NETWORK: "network",
  UNKNOWN: "unknown",
} as const;

export type ApiReason = (typeof API_REASON)[keyof typeof API_REASON];

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; reason: ApiReason };

export async function newUser(
  DisplayName: string,
): Promise<ApiResult<UserAPIData>> {
  try {
    const res = await nodeApi.newUser({ DisplayName });
    console.log(res);

    if (res.status === 200 && res.data) {
      const data: UserAPIData = {
        id: res.data.Profile.ID,
        name: res.data.Profile.Name,
        displayName: res.data.Profile.DisplayName,
      };
      return { ok: true, data };
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
  } catch (e) {
    console.log(e);
    return { ok: false, reason: API_REASON.NETWORK };
  }
}

export async function getUser(userID: UserID): Promise<ApiResult<UserAPIData>> {
  try {
    const res = await nodeApi.getUser(userID.id, userID.name);
    if (res.status === 200 && res.data) {
      const data: UserAPIData = {
        id: res.data.Profile.ID,
        name: res.data.Profile.Name,
        displayName: res.data.Profile.DisplayName,
      };
      return { ok: true, data };
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
  } catch (e) {
    return { ok: false, reason: API_REASON.NETWORK };
  }
}
