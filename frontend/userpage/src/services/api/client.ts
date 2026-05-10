import { handleAPI, type ApiResult } from "./handle-api";
import { setupClient } from "./client-setup";

import { newUser as _newUser, getUser as _getUser } from "./generated/sdk.gen";
import type { User } from "./generated/types.gen";

setupClient();

export async function newUser(displayName: string): Promise<ApiResult<User>> {
  return handleAPI(
    () => _newUser({ body: { DisplayName: displayName } }),
    (data) => data,
  );
}

export async function getUser(id: number, name: string): Promise<ApiResult<User>> {
  return handleAPI(
    () => _getUser({ path: { ID: id, Name: name } }),
    (data) => data,
  );
}
