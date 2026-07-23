import { handleAPI, type ApiResult } from "./handle-api";

import {
  auth as _auth,
  listUsers as _listUsers,
  enableUser as _enableUser,
  disableUser as _disableUser,
  deleteUser as _deleteUser,
  newUser as _newUser,
  listNodes as _listNodes,
  startNode as _startNode,
  stopNode as _stopNode,
  newNode as _newNode,
  deleteNode as _deleteNode,
  getSettings as _getSettings,
  setSettings as _setSettings,
} from "./generated/sdk.gen";
import type { User, Node, AuthResponse, UserView, Settings } from "./generated/types.gen";

export async function auth(pwd: string): Promise<ApiResult<AuthResponse>> {
  return handleAPI(
    () => _auth({ body: { password: pwd } }),
    (data) => data,
  );
}

export async function listUsers(): Promise<ApiResult<Array<UserView>>> {
  return handleAPI(
    () => _listUsers(),
    (data) => data.Users,
  );
}

export async function enableUser(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _enableUser({ body: { ID: id } }),
    () => {},
  );
}

export async function disableUser(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _disableUser({ body: { ID: id } }),
    () => {},
  );
}

export async function newUser(displayName: string): Promise<ApiResult<User>> {
  return handleAPI(
    () => _newUser({ body: { DisplayName: displayName } }),
    (data) => data,
  );
}

export async function deleteUser(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _deleteUser({ body: { ID: id } }),
    () => {},
  );
}

export async function listNodes(): Promise<ApiResult<Array<Node>>> {
  return handleAPI(
    () => _listNodes(),
    (data) => data.Nodes,
  );
}

export async function startNode(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _startNode({ body: { ID: id } }),
    () => {},
  );
}

export async function stopNode(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _stopNode({ body: { ID: id } }),
    () => {},
  );
}

export async function newNode(endpoint: string, accessKey: string): Promise<ApiResult<Node>> {
  return handleAPI(
    () => _newNode({ body: { Endpoint: endpoint, AccessKey: accessKey } }),
    (data) => data.Node,
  );
}

export async function deleteNode(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _deleteNode({ body: { ID: id } }),
    () => {},
  );
}

export async function getSettings(): Promise<ApiResult<Settings>> {
  return handleAPI(
    () => _getSettings(),
    (data) => data,
  );
}

export async function setSettings(settings: Settings): Promise<ApiResult<void>> {
  return handleAPI(
    () => _setSettings({ body: settings }),
    () => {},
  );
}
