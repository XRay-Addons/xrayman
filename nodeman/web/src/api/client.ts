import { getNodeManagementAPI } from "./generated";
import { UserID, UserAPIData } from "../lib/types";
import axios from "axios";

export const api = axios.create({
  baseURL: "http://localhost:80/api",
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
        subscriptionURL: "",
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

    //const res = await apiGetUser(userID.id, userID.name);

    if (res.status === 200 && res.data) {
      const data: UserAPIData = {
        id: res.data.Profile.ID,
        name: res.data.Profile.Name,
        displayName: res.data.Profile.DisplayName,
        subscriptionURL: "",
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
/*type User = components["schemas"]["User"];

const client = createClient<paths>({
  baseUrl: import.meta.env.DEV ? "http://localhost:8080" : "",
});

// Типы для ответов
export type User = {
  id: number;
  name: string;
  displayName: string;
  vlessUUID: string;
  targetStatus: string;
};

// Тип для результата API
export type ApiResult<T> =
  | { success: true; data: T }
  | {
      success: false;
      error: { status?: number; message: string; details?: string };
    };

// API методы
export const api = {

  createUser: async (displayName: string): Promise<ApiResult<User>> => {
    try {
      const { data, error } = await client.POST("/user/new", {
        body: { DisplayName: displayName },
      });

      if (error) {
        return {
          success: false,
          error: {
            status: error.status,
            message: error.Message || "Ошибка создания пользователя",
            details: error.Details,
          },
        };
      }

      if (!data) {
        return {
          success: false,
          error: { message: "Сервер вернул пустой ответ" },
        };
      }

      return {
        success: true,
        data: {
          id: data.Profile.ID,
          name: data.Profile.Name,
          displayName: data.Profile.DisplayName,
          vlessUUID: data.Profile.VlessUUID,
          targetStatus: data.TargetStatus,
        },
      };
    } catch (err) {
      return {
        success: false,
        error: {
          message: err instanceof Error ? err.message : "Неизвестная ошибка",
        },
      };
    }
  },


  fetchUser: async (id: number, name: string): Promise<ApiResult<User>> => {
    try {
      const { data, error } = await client.GET("/user/{ID}-{Name}", {
        params: {
          path: {
            ID: id,
            Name: name,
          },
        },
      });

      if (error) {
        return {
          success: false,
          error: {
            status: error.status,
            message: error.Message || "Пользователь не найден",
            details: error.Details,
          },
        };
      }

      if (!data) {
        return {
          success: false,
          error: { message: "Сервер вернул пустой ответ" },
        };
      }

      return {
        success: true,
        data: {
          id: data.Profile.ID,
          name: data.Profile.Name,
          displayName: data.Profile.DisplayName,
          vlessUUID: data.Profile.VlessUUID,
          targetStatus: data.TargetStatus,
        },
      };
    } catch (err) {
      return {
        success: false,
        error: {
          message: err instanceof Error ? err.message : "Неизвестная ошибка",
        },
      };
    }
  },
};

// Экспортируем отдельно для удобства
export const createUser = api.createUser;
export const fetchUser = api.fetchUser;*/
