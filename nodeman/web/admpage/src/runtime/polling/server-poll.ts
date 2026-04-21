import { reloadUsers } from "@/actions/users";
import { reloadNodes } from "@/actions/nodes";

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
        await Promise.all([
          reloadUsers({ quiet: true }),
          reloadNodes({ quiet: true }),
        ]);
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

const serverPoll = new ServerPoll();

export function startServerPoll() {
  serverPoll.start();
}

export function stopServerPoll() {
  serverPoll.stop();
}
