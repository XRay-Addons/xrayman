export function fixIOSActive() {
  document.addEventListener("touchstart", function () {}, { passive: true });
}
