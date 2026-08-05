import { getAuthToken } from "@/state/token";
import { client } from "./generated/client.gen";
import { authMan } from "./auth-man";
import { config } from "@/config/config";

let initialized = false;

export function ensureClient() {
  if (initialized) {
    return;
  }

  client.setConfig({
    auth: getToken,
    baseUrl: config().ApiPrefix,
    fetch: authFetch,
  });

  initialized = true;
}

function getToken(): string {
  return getAuthToken() ?? "";
}

async function authFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  const response = await fetch(input, init);
  if (!isAuthError(response)) {
    return response;
  }

  const requestBearer = extractTokenBearer(input, init);
  const actualBearer = makeTokenBearer(getAuthToken());
  if (requestBearer !== actualBearer) {
    return fetchWithTokenBearer(input, init, actualBearer);
  }

  const token = await authMan.waitForLogin();
  return fetchWithTokenBearer(input, init, makeTokenBearer(token));
}

function isAuthError(response: Response): boolean {
  return response.status === 401;
}

function extractTokenBearer(input: RequestInfo | URL, init?: RequestInit): string | null {
  const headers = new Headers(
    init?.headers ?? (input instanceof Request ? input.headers : undefined),
  );
  return headers.get("Authorization");
}

function makeTokenBearer(token: string | null): string | null {
  return token ? `Bearer ${token}` : null;
}

function fetchWithTokenBearer(
  input: RequestInfo | URL,
  init: RequestInit | undefined,
  bearer: string | null,
): Promise<Response> {
  const headers = new Headers(
    init?.headers ?? (input instanceof Request ? input.headers : undefined),
  );

  if (bearer) {
    headers.set("Authorization", bearer);
  }

  return fetch(input, {
    ...init,
    headers,
  });
}
