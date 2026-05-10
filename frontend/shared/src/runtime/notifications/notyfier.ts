import { Notyf } from "notyf";

export class Notyfier {
  private notyf?: Notyf;

  constructor() {
    this.notyf = new Notyf({
      duration: 6000000,
      position: { x: "right", y: "bottom" },
      dismissible: true,
    });
  }

  errorNotification(message: string, details?: string) {
    const error = details ? `<b>${message}</b><br/>${details}` : `<b>${message}</b>`;
    this.notyf!.error({ message: error });
    console.log("error:", message, details);
  }
}
