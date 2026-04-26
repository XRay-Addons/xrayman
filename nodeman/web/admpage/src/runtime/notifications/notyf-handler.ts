import { Notyf } from "notyf";
import { setErrorHandler } from "./errors";

const notyf = new Notyf({
  duration: 60000,
  position: { x: "right", y: "bottom" },
  dismissible: true,
});

export function notyfErrorHandler(message: string, description?: string) {
  const error = `<b>${message}</b><br/>${description}`;
  notyf.error({ message: error });
}
