import { getAuthToken } from "@/state/token";
import { client } from "./generated/client.gen";
import { authMan } from "./auth-man";
import { config } from "@/config/config";

export function setupClient() {
  client.setConfig({
    auth: getToken,
    baseUrl: config.ApiPrefix,
    fetch: authFetch,
  });
}

function getToken(): string {
  return getAuthToken() ?? "[no token]";
}

async function authFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  let response = await fetch(input, init);
  while (isJwtIssue(input, response)) {
    console.log("auth JWT issue");
    await authMan.handle401(async () => {
      console.log("auth 401, fetch again");
      response = await fetch(input, withHeader(init, "Authorization", `Bearer ${getToken()}`));
    });
  }
  return response;
}

function isJwtIssue(request: RequestInfo | URL, response: Response): boolean {
  return (
    response.status == 401 && request instanceof Request && request.headers.has("Authorization")
  );
}

function withHeader(init: RequestInit | undefined, key: string, value: string): RequestInit {
  const headers = new Headers(init?.headers);

  headers.set(key, value);

  return {
    ...init,
    headers,
  };
}
