export const API_REASON = {
  BAD_REQUEST: "bad_request",
  NOT_FOUND: "not_found",
  UNAUTHORIZED: "unauthorized",
  FORBIDDEN: "forbidden",
  NETWORK: "network",
  UNKNOWN: "unknown",
} as const;

export type ApiReason = (typeof API_REASON)[keyof typeof API_REASON];

export interface Column {
  label: string;
  path: string;
  format?: (value: any, row: any) => string;
}
