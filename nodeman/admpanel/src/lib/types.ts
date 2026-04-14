export type UserID = {
  id: number;
  name: string;
};

export type UserAPIData = {
  id: number;
  name: string;
  displayName: string;
};

export type User = {
  id: number;
  name: string;
  displayName: string;
};

export const API_REASON = {
  BAD_REQUEST: "bad_request",
  NOT_FOUND: "not_found",
  UNAUTHORIZED: "unauthorized",
  FORBIDDEN: "forbidden",
  NETWORK: "network",
  UNKNOWN: "unknown",
} as const;

export type ApiReason = (typeof API_REASON)[keyof typeof API_REASON];
