function openShareWindow(url: string) {
  const width = 480;
  const height = 320;

  const left = (window.screen.width - width) / 2;
  const top = (window.screen.height - height) / 2;

  const features = [
    "toolbar=no",
    "location=no",
    "status=no",
    "menubar=no",
    "scrollbars=yes",
    "resizable=yes",
    `width=${width}`,
    `height=${height}`,
    `top=${top}`,
    `left=${left}`,
  ].join(",");

  window.open(url, "Share", features);
}

export function ShareViaTG(url: string, text: string) {
  const tgLink = `//telegram.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(text)}`;
  openShareWindow(tgLink);
}
