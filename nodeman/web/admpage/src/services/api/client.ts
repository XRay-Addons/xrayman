import { API_REASON, type ApiReason } from "./api-reason";
import { config } from "@/config/config";
import { client } from "./generated/client.gen";
import {
  listUsers as _listUsers,
  enableUser as _enableUser,
  disableUser as _disableUser,
  newUser as _newUser,
  listNodes as _listNodes,
  startNode as _startNode,
  stopNode as _stopNode,
  newNode as _newNode,
} from "./generated/sdk.gen";
import type { Error, User, Node } from "./generated/types.gen";

client.setConfig({
  //baseUrl: config.API_URLPATH,
  baseUrl: "http://localhost:80/api",
});

export type ApiResult<T> =
  | { ok: true; data: T }
  | { ok: false; reason: ApiReason };

type ApiResponse<T> = (
  | { data: T; error: undefined }
  | { data: undefined; error: Error }
) & {
  request: Request;
  response: Response;
};

export async function listUsers(): Promise<ApiResult<Array<User>>> {
  return handleAPI(
    () => _listUsers(),
    (data) => data.Users,
  );
}

export async function enableUser(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _enableUser({ body: { ID: id } }),
    (data) => {},
  );
}

export async function disableUser(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _disableUser({ body: { ID: id } }),
    (data) => {},
  );
}

export async function newUser(displayName: string): Promise<ApiResult<User>> {
  return handleAPI(
    () => _newUser({ body: { DisplayName: displayName } }),
    (data) => data,
  );
}

export async function listNodes(): Promise<ApiResult<Array<APINode>>> {
  return handleAPI(
    () => _listNodes(),
    (data) => data.Nodes,
  );
}

export async function startNode(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _startNode({ body: { ID: id } }),
    (data) => {},
  );
}

export async function stopNode(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _stopNode({ body: { ID: id } }),
    (data) => {},
  );
}

export async function newNode(
  endpoint: string,
  accessKey: string,
): Promise<ApiResult<Node>> {
  return handleAPI(
    () => _newNode({ body: { Endpoint: endpoint, AccessKey: accessKey } }),
    (data) => data,
  );
}

async function handleAPI<T, R>(
  apiCall: () => Promise<ApiResponse<T>>,
  transform: (data: T) => R,
): Promise<ApiResult<R>> {
  try {
    console.log("call api");
    const resp = await apiCall();
    console.log("call api response:", resp);

    if (!resp.error) {
      return {
        ok: true,
        data: transform(resp.data),
      };
    }
    console.log("api call error:", resp.error);

    const status = resp.response.status;
    let reason: ApiReason;

    switch (status) {
      case 400:
        reason = API_REASON.BAD_REQUEST;
        break;
      case 401:
        reason = API_REASON.UNAUTHORIZED;
        break;
      case 403:
        reason = API_REASON.FORBIDDEN;
        break;
      case 404:
        reason = API_REASON.NOT_FOUND;
        break;
      default:
        reason = API_REASON.UNKNOWN;
    }

    return { ok: false, reason };
  } catch (error) {
    console.log("api call error:", error);
    return { ok: false, reason: API_REASON.NETWORK };
  }
}
