import { User, fetchUser } from "./api";
import { HAPP_INTENT } from "../config/config.json";

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

function setState(state: STATE, root: HTMLElement = document.body) {
  const visibleSelectors = STATE_MAP[state] ?? [];

  for (const id of ALL_BLOCKS) {
    const el = document.getElementById(id);
    if (el) el.style.display = "none";
  }

  for (const id of visibleSelectors) {
    const el = document.getElementById(id);
    if (el) el.style.display = "";
  }
  const body = document.getElementById("body");
  if (body) body.classList.remove("invisible");
}
setState(STATE.LOGGED);

function setLoggedState(user: User) {
  // title
  const titleEl = document.getElementById("title-name");
  if (titleEl) titleEl.innerHTML = user.visibleName;

  // input
  const inputEl = document.getElementById("user-url-input") as HTMLInputElement;
  if (inputEl) inputEl.value = window.location.href;

  // happ btn
  const happURL = HAPP_INTENT + user.subscriptionURL;
  const btnEl = document
    .getElementById("run-app-btn")
    ?.querySelector<HTMLButtonElement>("button");
  if (btnEl) {
    btnEl.addEventListener("click", () => {
      window.open(happURL, "_blank", "noopener");
    });
  }

  setState(STATE.LOGGED);
}

function setUnloggedState() {
  const titleEl = document.getElementById("title-name");
  if (titleEl) {
    const placeholder = titleEl.getAttribute("placeholder") as string;
    titleEl.innerHTML = placeholder;
  }

  const inputEl = document.getElementById("user-url-input") as HTMLInputElement;
  if (inputEl) inputEl.value = "";

  const btnEl = document.getElementById("run-app-btn");
  console.log(btnEl);
  if (btnEl) btnEl.setAttribute("data-url", "");

  setState(STATE.UNLOGGED);
}

function setCookie(name: string, value: string) {
  localStorage.setItem(name, value);
}

function getCookie(name: string): string | null {
  return localStorage.getItem(name);
}

function getUserFromPath() {
  const match = window.location.pathname.match(/^\/(\d+)-([^/]+)$/);
  return match ? { id: match[1], name: match[2] } : null;
}

function setUserToPath(u) {
  history.replaceState({}, "", `/${u.id}-${u.name}`);
}

(async function init() {
  const userProps = getUserFromPath();
  let id = userProps?.id || getCookie("UserID");
  let name = userProps?.name || getCookie("Name");

  if (id && name) {
    try {
      const user = await fetchUser(id, name);
      if (user) {
        setCookie("UserID", id);
        setCookie("Name", name);
        setUserToPath(userProps);
        setLoggedState(user);
        return;
      }
    } catch (e) {
      console.error(e);
    }
  }
  setUnloggedState();
})();

// setUnloggedState();

const user: User = {
  id: "16-stepan",
  visibleName: "ИГОРЬ",
  subscriptionURL: "https://happ.su/subscription.json",
};

setLoggedState(user);
