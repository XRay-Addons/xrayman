export function upgradeProperties(el: HTMLElement, props: readonly string[]) {
  for (const prop of props) {
    if (Object.prototype.hasOwnProperty.call(el, prop)) {
      const value = (el as any)[prop];
      delete (el as any)[prop];
      (el as any)[prop] = value;
      console.log("property ${prop} updated");
    } else {
      console.log("property ${prop} not exists yet");
    }
  }
}
