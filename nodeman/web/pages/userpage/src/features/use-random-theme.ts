import { setPalette } from "@/actions/palette";
import palettes from "@/data/palettes/palettes.json";
import { RandomTheme, type Palette } from "@xrayman/shared/features/random-theme";

var randomTheme: RandomTheme;

export function useRandomTheme() {
  console.log("user random theme");
  randomTheme = new RandomTheme(palettes, setPalette);
}
