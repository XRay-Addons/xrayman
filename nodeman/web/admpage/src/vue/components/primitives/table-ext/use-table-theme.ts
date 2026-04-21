import { type Palette, PaletteItem } from "@/state/palette";
import { type ShallowRef, computed } from "vue";
import Color from "colorjs.io";

export function useTableTheme(paletteRef: ShallowRef<Palette>) {
  return computed(() => {
    const mainColor = new Color(paletteRef.value[PaletteItem.Table]);
    const textColor = contrastedColor(mainColor);

    return {
      token: {
        colorBgContainer: mainColor.toString(),
        colorTextHeading: textColor.toString(),
        colorText: textColor.toString(),
      },
      components: {
        Table: {
          borderRadiusSM: 0,
          borderRadiusMD: 0,
          borderRadiusLG: 0,
          algorithm: true,
        },
      },
    };
  });
}

function contrastedColor(color: Color): Color {
  const b = new Color("black");
  const w = new Color("white");

  const contrastB = color.contrast(new Color("black"), "WCAG21");
  const contrastW = color.contrast(new Color("white"), "WCAG21");
  const contrastedColor = contrastB > contrastW ? b : w;

  return color.mix(contrastedColor, 0.75, { space: "srgb" });
}
