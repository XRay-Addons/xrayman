import { ref, shallowRef, type ShallowRef } from "vue";

export enum PaletteItem {
  BG = "BG",
  Card = "Card",
  Title = "Title",
  Button = "Button",
  Input = "Input",
  Table = "Table",
}

let paletteValue = {
  [PaletteItem.BG]: "#feffa3",
  [PaletteItem.Card]: "#d4ffea",
  [PaletteItem.Title]: "#ffd4e5",
  [PaletteItem.Button]: "#eecbff",
  [PaletteItem.Input]: "#dbdcff",
  [PaletteItem.Table]: "#d4ffea",
};

const palette = shallowRef(paletteValue);

export type Palette = typeof paletteValue;

export function getPaletteState() {
  return palette.value;
}

export function setPaletteState(p: Palette) {
  palette.value = p;
}

export function getPaletteRef(): ShallowRef<Palette> {
  return palette;
}
