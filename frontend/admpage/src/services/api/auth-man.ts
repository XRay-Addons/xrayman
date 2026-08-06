import { getAuthToken, setAuthToken } from "@/state/token";

type LoginPromise = Promise<string>;

export const loginRequiredEvent = "auth:login-required";

class AuthMan {
  private loginPromise: LoginPromise | null = null;
  private resolveLogin: ((token: string) => void) | null = null;
  private rejectLogin: ((error: unknown) => void) | null = null;

  waitForLogin(): Promise<string> {
    if (this.loginPromise) {
      return this.loginPromise;
    }

    this.loginPromise = new Promise<string>((resolve, reject) => {
      this.resolveLogin = resolve;
      this.rejectLogin = reject;

      window.dispatchEvent(new CustomEvent(loginRequiredEvent));
    }).finally(() => {
      this.loginPromise = null;
      this.resolveLogin = null;
      this.rejectLogin = null;
    });

    return this.loginPromise;
  }

  onLoginSuccess(token: string) {
    setAuthToken(token);

    this.resolveLogin?.(token);
  }

  onLoginFail(error?: unknown) {
    this.rejectLogin?.(error ?? new Error("auth failed"));
  }

  getToken() {
    return getAuthToken();
  }
}

export const authMan = new AuthMan();
