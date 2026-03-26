import { User } from "./types";

class UserStore {
  private user: User | null = null;
  private listeners = new Set<() => void>();

  get() {
    return this.user;
  }

  set(user: User | null) {
    this.user = user;
    this.listeners.forEach((l) => l());
  }

  subscribe(fn: () => void) {
    this.listeners.add(fn);
    return () => this.listeners.delete(fn);
  }
}

export const userStore = new UserStore();
//window.userStore = userStore;
