<script setup lang="ts">
import { watch, onMounted } from "vue";
import { colors, Colors } from "../../state/colors";

// map CSS <-> state
const map = {
  [Colors.BG]: "--bg-color",
  [Colors.Card]: "--card-color",
  [Colors.Title]: "--title-color",
  [Colors.Button]: "--button-color",
  [Colors.Input]: "--input-color",
  [Colors.Table]: "--table-color",
} as const;

// sync CSS → state
function colorsFromCSS(root: HTMLElement) {
  const styles = getComputedStyle(root);

  for (const key in map) {
    const cssVar = map[key as Colors];
    const value = styles.getPropertyValue(cssVar).trim();

    if (value) {
      colors.value[key as Colors] = value;
    }
  }
}

// sync: state → CSS
function colorsToCSS(root: HTMLElement) {
  for (const key in map) {
    const cssVar = map[key as Colors];
    root.style.setProperty(cssVar, colors.value[key as Colors]);
  }
}

// css values are source of initial values.
// reactive color updates leads to css updates
function getRoot() {
  return document.documentElement;
}

function initGlobals() {
  const root = getRoot();
  if (!root) return;
  colorsFromCSS(root);
}

function applyGlobals() {
  const root = getRoot();
  if (!root) return;
  colorsToCSS(root);
}

watch(
  colors,
  () => {
    applyGlobals();
  },
  { deep: true },
);

onMounted(() => {
  initGlobals();
});
</script>

<template>
  <slot />
</template>
