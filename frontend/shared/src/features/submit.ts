export function addSubmitListener(
  btnEl: HTMLButtonElement,
  formEl: HTMLFormElement,
  submitFn: () => Promise<void>,
): void {
  btnEl.addEventListener("click", (e) => {
    e.preventDefault();
    formEl.requestSubmit();
  });

  formEl.addEventListener("submit", async (e) => {
    e.preventDefault();

    if (btnEl.disabled) {
      return;
    }

    btnEl.disabled = true;

    try {
      await submitFn();
    } finally {
      btnEl.disabled = false;
    }
  });
}
