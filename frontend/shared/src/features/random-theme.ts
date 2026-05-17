// change theme on double tap, set it to CSS and additionaly via setPalette

export enum PaletteItem {
  BG = "BG",
  Card = "Card",
  Title = "Title",
  Button = "Button",
  Input = "Input",
  Table = "Table",
  Success = "Success",
  Tag = "Tag",
}

export type Palette = Record<PaletteItem, string>;

export type OnSetPalette = (p: Palette) => void;

export class RandomTheme {
  constructor(palettes: string[][], callback: OnSetPalette) {
    let lastTap = 0;
    const timeout = 300;

    // init palette by CSS
    const domPalette = getDOMPalette();
    callback(domPalette);

    document.addEventListener("click", () => {
      const now = Date.now();
      if (now - lastTap < timeout) {
        const next = getRandomPalette(palettes);
        if (next) {
          callback(next);
        }
      }

      lastTap = now;
    });
  }
}

function getRandomPalette(palettes: string[][]): Palette | null {
  const p = palettes[Math.floor(Math.random() * palettes.length)].slice();
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
    [PaletteItem.Success]: p[0 % n], // success color = card color
    [PaletteItem.Tag]: p[0 % n], // tag color = card color
  };
}

function getDOMPalette(): Palette {
  const root = document.documentElement;
  return {
    [PaletteItem.BG]: getComputedStyle(root).getPropertyValue("--bg-color"),
    [PaletteItem.Card]: getComputedStyle(root).getPropertyValue("-card-color"),
    [PaletteItem.Title]: getComputedStyle(root).getPropertyValue("--title-color"),
    [PaletteItem.Button]: getComputedStyle(root).getPropertyValue("--button-color"),
    [PaletteItem.Input]: getComputedStyle(root).getPropertyValue("--input-color"),
    [PaletteItem.Table]: getComputedStyle(root).getPropertyValue("--table-color"),
    [PaletteItem.Success]: getComputedStyle(root).getPropertyValue("--success-color"),
    [PaletteItem.Tag]: getComputedStyle(root).getPropertyValue("--tag-color"),
  };
}

export function setDOMPalette(p: Palette) {
  const root = document.documentElement;
  root.style.setProperty("--bg-color", p[PaletteItem.BG]);
  root.style.setProperty("--card-color", p[PaletteItem.Card]);
  root.style.setProperty("--title-color", p[PaletteItem.Title]);
  root.style.setProperty("--button-color", p[PaletteItem.Button]);
  root.style.setProperty("--input-color", p[PaletteItem.Input]);
  root.style.setProperty("--table-color", p[PaletteItem.Table]);
  root.style.setProperty("--success-color", p[PaletteItem.Success]);
  root.style.setProperty("--tag-color", p[PaletteItem.Tag]);
}
