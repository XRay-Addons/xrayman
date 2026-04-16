import { syncUsers } from "@/state/users";
import { syncNodes } from "@/state/nodes";

class ServerPoll {
  private timer: number | null = null;
  private running = false;
  private interval = 3200;

  start() {
    if (this.running) return;
    this.running = true;

    const tick = async () => {
      if (!this.running) return;

      try {
        await Promise.all([syncUsers(), syncNodes()]);
      } catch (e) {
        console.error("[server-poll] sync error", e);
      } finally {
        if (this.running) {
          this.timer = window.setTimeout(tick, this.interval);
        }
      }
    };

    tick();
  }

  stop() {
    this.running = false;

    if (this.timer) {
      clearTimeout(this.timer);
      this.timer = null;
    }
  }
}

export const serverPoll = new ServerPoll();
