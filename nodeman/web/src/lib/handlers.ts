export enum STATE {
  UNLOGGED = "unlogged",
  LOGGED = "logged",
}

const STATE_MAP: Record<STATE, string[]> = {
  [STATE.UNLOGGED]: ["title", "choose-name-card"],
  [STATE.LOGGED]: [
    "title",
    "user-url-card",
    "install-app-card",
    "run-app-card",
  ],
};
const ALL_BLOCKS = Object.values(STATE_MAP).flat();

export function setState(state: STATE, root: HTMLElement = document.body) {
  const visibleSelectors = STATE_MAP[state] ?? [];

  for (const id of ALL_BLOCKS) {
    const el = document.getElementById(id);
    if (el) el.style.display = "none";
  }

  for (const id of visibleSelectors) {
    const el = document.getElementById(id);
    if (el) el.style.display = "";
  }
}
