export const State = {
  Idle: "idle",
  LoggedIn: "logged-in",
  LoggedOut: "logged-out",
  ServerError: "server-error",
} as const;

export type State = (typeof State)[keyof typeof State];
