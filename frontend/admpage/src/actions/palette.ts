import { setPaletteState } from "@/state/palette";
import { type Palette, setDOMPalette } from "@xrayman/shared/features/random-theme";

export function setPalette(p: Palette) {
  setPaletteState(p);
  setDOMPalette(p);
}
