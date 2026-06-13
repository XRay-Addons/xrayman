export function createPoll(pollFn: () => Promise<void>, interval = 32000) {
  let stopped = false;
  let timer: number | null = null;

  const loop = async () => {
    if (stopped) return;

    try {
      await pollFn();
    } catch (e) {
      console.error("[poll] error", e);
    } finally {
      if (!stopped) {
        timer = window.setTimeout(loop, interval);
      }
    }
  };

  return {
    start() {
      if (timer !== null) return;
      timer = window.setTimeout(loop, interval);
    },
    stop() {
      stopped = true;
      if (timer !== null) {
        clearTimeout(timer);
        timer = null;
      }
    },
  };
}
