import { setPalette } from "@/actions/palette";
import palettes from "@/data/palettes/palettes.json";
import { PaletteItem, type Palette } from "@/state/palette";

export function initRandomThemeFeature() {
  let lastTap = 0;
  const timeout = 300;

  document.addEventListener("click", () => {
    const now = Date.now();

    if (now - lastTap < timeout) {
      const next = getRandomPalette();
      if (next) setPalette(next);
    }

    lastTap = now;
  });
}

function getRandomPalette(): Palette | null {
  const p = palettes[Math.floor(Math.random() * palettes.length)]?.slice();
  if (!p?.length) return null;

  p.sort(() => Math.random() - 0.5);
  const n = p.length;

  return {
    [PaletteItem.BG]: p[0 % n],
    [PaletteItem.Card]: p[1 % n],
    [PaletteItem.Title]: p[2 % n],
    [PaletteItem.Button]: p[3 % n],
    [PaletteItem.Input]: p[4 % n],
    [PaletteItem.Table]: p[1 % n], // table color = card color
  };
}
