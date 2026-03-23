export const UsernameCreate = "username-create";

export interface UsernameCreateDetail {
  displayName: string;
}

export function emitUsernameCreate(
  host: HTMLElement,
  detail: UsernameCreateDetail,
) {
  const event = new CustomEvent<UsernameCreateDetail>(UsernameCreate, {
    bubbles: true,
    composed: true,
    detail,
  });
  host.dispatchEvent(event);
}

export function onUsernameCreate(
  host: HTMLElement,
  callback: (event: CustomEvent<UsernameCreateDetail>) => void,
) {
  host.addEventListener(UsernameCreate, callback as EventListener);
}
