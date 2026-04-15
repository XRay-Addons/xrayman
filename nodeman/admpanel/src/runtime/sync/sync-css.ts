import { watch, onMounted } from "vue";
import { colors, Colors, type ColorsType } from "../../state/colors";

// map CSS <-> state
const map = {
  [Colors.BG]: "--bg-color",
  [Colors.Card]: "--card-color",
  [Colors.Title]: "--title-color",
  [Colors.Button]: "--button-color",
  [Colors.Input]: "--input-color",
  [Colors.Table]: "--table-color",
} as const;

// sync CSS -> state
function colorsFromCSS(root: HTMLElement) {
  // read colors from css, update by single op.
  const updatedColors = { ...colors.value };
  const styles = getComputedStyle(root);
  Object.entries(map).forEach(([key, cssVar]) => {
    const value = styles.getPropertyValue(cssVar).trim();
    if (value) updatedColors[key as Colors] = value;
  });
  colors.value = updatedColors;
}

// sync: state -> CSS
function colorsToCSS(colors: ColorsType, root: HTMLElement) {
  Object.entries(map).forEach(([key, cssVar]) => {
    root.style.setProperty(cssVar, colors[key as Colors]);
  });
}

// css values are source of initial values.
// reactive color updates leads to css updates
export function initFromCSS() {
  const root = document.documentElement;
  if (!root) return;
  colorsFromCSS(root);
}

function applyToCSS(colors) {
  const root = document.documentElement;
  if (!root) return;
  colorsToCSS(colors, root);
}

watch(
  colors,
  (newColors) => {
    applyToCSS(newColors);
  },
  { deep: true },
);
