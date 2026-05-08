export const API_REASON = {
  BAD_REQUEST: "bad_request",
  NOT_FOUND: "not_found",
  UNAUTHORIZED: "unauthorized",
  FORBIDDEN: "forbidden",
  NETWORK: "network",
  UNKNOWN: "unknown",
} as const;

export type ApiReason = (typeof API_REASON)[keyof typeof API_REASON];

export function getErrorReason(status: number): ApiReason {
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
