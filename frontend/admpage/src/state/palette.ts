import { shallowRef, type ShallowRef } from "vue";
import { type Palette, PaletteItem } from "@xrayman/shared/features/random-theme";

let paletteValue: Palette = {
  [PaletteItem.BG]: "transparent",
  [PaletteItem.Card]: "transparent",
  [PaletteItem.Title]: "transparent",
  [PaletteItem.Button]: "transparent",
  [PaletteItem.Input]: "transparent",
  [PaletteItem.Table]: "transparent",
  [PaletteItem.Success]: "transparent",
};

const palette = shallowRef(paletteValue);

export function getPaletteState() {
  return palette.value;
}

export function setPaletteState(p: Palette) {
  palette.value = p;
}

export function getPaletteRef(): ShallowRef<Palette> {
  return palette;
}
