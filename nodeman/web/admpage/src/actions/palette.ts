import { setPaletteState, type Palette } from "@/state/palette";
import { updateDOMPalette } from "@/runtime/dom/palette-dom";

export function setPalette(p: Palette) {
  setPaletteState(p);
  updateDOMPalette();
}
