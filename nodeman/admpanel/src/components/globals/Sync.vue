<script setup lang="ts">
import { watch, onMounted, inject } from "vue";
import { colors, Colors } from "./colors";

/* ========================================
  global colors. css is initial data source
  ====================================== */
const map = {
  [Colors.BG]: "--bg-color",
  [Colors.Card]: "--card-color",
  [Colors.Title]: "--title-color",
  [Colors.Button]: "--button-color",
  [Colors.Input]: "--input-color",
  [Colors.Table]: "--table-color",
} as const;

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

function colorsToCSS(root: HTMLElement) {
  for (const key in map) {
    const cssVar = map[key as Colors];
    root.style.setProperty(cssVar, colors.value[key as Colors]);
  }
}

/* ========================================
  init globals by CSS, apply globals changes to CSS
  ====================================== */
function initGlobals() {
  const root = document.documentElement;
  if (!root) {
    console.warn("globals init from CSS fiasco");
  } // too early, init order fiasco
  colorsFromCSS(root);
}

function applyGlobals() {
  const root = document.documentElement;
  if (!root) {
    return;
  } // too early, init order fiasco
  colorsToCSS(root);
}

watch(
  colors,
  () => {
    applyGlobals();
  },
  { deep: true, immediate: false },
);

onMounted(() => {
  initGlobals();
});
</script>

<template>
  <slot />
</template>
