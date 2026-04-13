<script setup lang="ts">
import { onMounted, onUnmounted, ref, provide } from "vue";
import emitter, { Events, Colors } from "./Globals.ts";
import type { SetColorDetail } from "./Globals.ts";
import GlobalsToCSS from "./GlobalsToCSS.vue";

const colors = ref({
  [Colors.Background]: "#ff00ff",
  [Colors.Foreground]: "#ffffff",
});

function handleSetColor(payload: SetColorDetail) {
  console.log("received set color", payload);
  colors.value[payload.color] = payload.value;
}

onMounted(() => {
  emitter.on(Events.SetColor, handleSetColor);
});

onUnmounted(() => {
  emitter.off(Events.SetColor, handleSetColor);
});

provide("colors", colors);

console.log("Globals: colors provided", colors.value);
</script>

<template>
  <GlobalsToCSS>
    <slot />
  </GlobalsToCSS>
</template>
