<script setup lang="ts">
import { watch, onMounted, inject } from "vue";
import { Colors } from "./Globals.ts";
import type { SetColorDetail } from "./Globals.ts";

const colors = inject("colors");

function colorsToCSS(root: HTMLRootElement) {
  root.style.setProperty("--bg-color", colors.value[Colors.Background]);
  root.style.setProperty("--foreground-color", colors.value[Colors.Foreground]);
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
