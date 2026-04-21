import { getPaletteState, type Palette, PaletteItem } from "@/state/palette";

export function getDOMPalette(): Palette {
  const root = document.documentElement;
  return {
    [PaletteItem.BG]: root.style.getPropertyValue("--bg-color"),
    [PaletteItem.Card]: root.style.getPropertyValue("-card-color"),
    [PaletteItem.Title]: root.style.getPropertyValue("--title-color"),
    [PaletteItem.Button]: root.style.getPropertyValue("--button-color"),
    [PaletteItem.Input]: root.style.getPropertyValue("--input-color"),
    [PaletteItem.Table]: root.style.getPropertyValue("--table-color"),
  };
}

export function updateDOMPalette() {
  const p = getPaletteState();
  const root = document.documentElement;
  root.style.setProperty("--bg-color", p[PaletteItem.BG]);
  root.style.setProperty("--card-color", p[PaletteItem.Card]);
  root.style.setProperty("--title-color", p[PaletteItem.Title]);
  root.style.setProperty("--button-color", p[PaletteItem.Button]);
  root.style.setProperty("--input-color", p[PaletteItem.Input]);
  root.style.setProperty("--table-color", p[PaletteItem.Table]);
}
