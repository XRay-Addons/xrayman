import { Notyf } from "notyf";

export class Notyfier {
  private notyf?: Notyf;

  constructor() {
    this.notyf = new Notyf({
      duration: 5000,
      position: { x: "right", y: "bottom" },
      dismissible: true,
    });
  }

  errorNotification(message: string, details?: string) {
    this._prepareContainer();
    const error = details ? `<b>${message}</b><br/>${details}` : `<b>${message}</b>`;
    this.notyf!.error({ message: error });
    console.log("error:", message, details);
  }

  _prepareContainer() {
    const notyfContainer = document.querySelector(".notyf");
    if (!notyfContainer) return;

    // attach notyf container to open modal dialog if opened
    const openDialogs = Array.from(document.querySelectorAll("dialog[open]"));
    if (openDialogs.length > 0) {
      const activeDialog = openDialogs[openDialogs.length - 1];
      if (notyfContainer.parentElement !== activeDialog) {
        activeDialog.appendChild(notyfContainer);
      }
    } else {
      if (notyfContainer.parentElement !== document.body) {
        document.body.appendChild(notyfContainer);
      }
    }
  }
}
