import {
  i18nateDOM,
  startI18nObserver,
  stopI18nObserver,
} from "./dom/i18nate-dom";
import { setErrorHandler } from "./notifications/errors";
import { notyfErrorHandler } from "./notifications/notyf-handler";

setErrorHandler(notyfErrorHandler);
i18nateDOM();
startI18nObserver();

window.addEventListener("beforeunload", () => {
  stopI18nObserver();
  setErrorHandler(null);
});
