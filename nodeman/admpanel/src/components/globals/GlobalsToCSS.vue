<script setup lang="ts">
import { watch, onMounted, inject } from "vue";
import { colors, Colors } from "./colors";

function colorsToCSS(root: HTMLElement) {
  root.style.setProperty("--bg-color", colors.value[Colors.BG]);
  root.style.setProperty("--card-color", colors.value[Colors.Card]);
  root.style.setProperty("--title-color", colors.value[Colors.Title]);
  root.style.setProperty("--button-color", colors.value[Colors.Button]);
  root.style.setProperty("--input-color", colors.value[Colors.Input]);
  root.style.setProperty("--table-color", colors.value[Colors.Table]);
}

function globalsToCSS() {
  const root = document.documentElement;
  if (!root) {
    return;
  } // too early, init order fiasco
  colorsToCSS(root);
}

watch(
  colors,
  () => {
    globalsToCSS();
  },
  { deep: true, immediate: false },
);

onMounted(() => {
  globalsToCSS();
});
</script>

<template>
  <slot />
</template>
