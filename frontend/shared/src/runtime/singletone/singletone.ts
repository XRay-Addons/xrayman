export function makeSingleton<T>(initFn: () => Promise<T>) {
  let value: T | null = null;
  let initPromise: Promise<T> | null = null;

  return {
    async get(): Promise<T> {
      if (value !== null) {
        return value;
      }

      if (!initPromise) {
        initPromise = initFn().then((v) => {
          value = v;
          return v;
        });
      }

      return initPromise;
    },

    reset() {
      value = null;
      initPromise = null;
    },
  };
}
