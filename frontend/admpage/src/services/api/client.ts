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
  listSubHeaders as _listSubHeaders,
  deleteSubHeader as _deleteSubHeader,
  newSubHeader as _newSubHeader,
} from "./generated/sdk.gen";
import type { User, Node, Header, AuthResponse } from "./generated/types.gen";

export async function auth(pwd: string): Promise<ApiResult<AuthResponse>> {
  return handleAPI(
    () => _auth({ body: { password: pwd } }),
    (data) => data,
  );
}

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

export async function deleteUser(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _deleteUser({ body: { ID: id } }),
    (data) => {},
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
    (data) => {},
  );
}

export async function stopNode(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _stopNode({ body: { ID: id } }),
    (data) => {},
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
    (data) => {},
  );
}

export async function listSubHeaders(): Promise<ApiResult<Array<Header>>> {
  return handleAPI(
    () => _listSubHeaders(),
    (data) => data.Headers,
  );
}

export async function newSubHeader(key: string, value: string): Promise<ApiResult<Header>> {
  return handleAPI(
    () => _newSubHeader({ body: { Key: key, Value: value } }),
    (data) => data,
  );
}

export async function deleteSubHeader(id: number): Promise<ApiResult<void>> {
  return handleAPI(
    () => _deleteSubHeader({ body: { ID: id } }),
    (data) => {},
  );
}
