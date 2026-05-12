export const config = {
  ApiPrefix: (window as any).__CONFIG__?.api_prefix ?? "/api",
  UserPrefix: (window as any).__CONFIG__?.user_prefix ?? "/u",
};
