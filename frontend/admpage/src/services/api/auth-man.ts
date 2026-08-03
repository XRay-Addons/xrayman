import { setAuthToken } from "@/state/token";

type RetryFn<T> = () => Promise<T>;

type PendingRequest = {
  resolve: (v: unknown) => void;
  reject: (e: unknown) => void;
  retry: RetryFn<unknown>;
};

export const loginRequiredEvent = "auth:login-required";

class AuthMan {
  private loginInProgress = false;
  private queue: PendingRequest[] = [];

  handle401<T>(retry: RetryFn<T>): Promise<T> {
    return new Promise<T>((resolve, reject) => {
      this.queue.push({
        resolve: resolve as (v: unknown) => void,
        reject,
        retry: retry as RetryFn<unknown>,
      });

      if (!this.loginInProgress) {
        this.startLoginFlow();
      }
    });
  }

  private startLoginFlow() {
    this.loginInProgress = true;

    window.dispatchEvent(new CustomEvent(loginRequiredEvent));
  }

  async onLoginSuccess(token: string) {
    setAuthToken(token);

    const queue = [...this.queue];
    this.queue = [];

    this.loginInProgress = false;

    for (const req of queue) {
      try {
        const result = await req.retry();
        req.resolve(result);
      } catch (e) {
        req.reject(e);
      }
    }
  }

  onLoginFail(error?: unknown) {
    const queue = [...this.queue];
    this.queue = [];

    this.loginInProgress = false;

    for (const req of queue) {
      req.reject(error ?? new Error("auth failed"));
    }
  }
}

export const authMan = new AuthMan();
