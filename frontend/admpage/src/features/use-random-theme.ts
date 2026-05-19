import { setPalette } from "@/actions/palette";
import palettes from "@/data/palettes/palettes.json";
import { UseRandomTheme } from "@xrayman/shared/features/random-theme";

export function useRandomTheme() {
  UseRandomTheme(palettes, setPalette);
}
