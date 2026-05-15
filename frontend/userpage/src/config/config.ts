export const config = {
  ApiPrefix: (window as any).__CONFIG__?.api_prefix ?? "http://localhost:1001/api",
  UserPagePrefix: (window as any).__CONFIG__?.user_prefix ?? "/",
  HAPP_INTENT: "happ://add/",
};
