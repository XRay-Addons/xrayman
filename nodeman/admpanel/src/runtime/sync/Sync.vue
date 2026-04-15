<script setup lang="ts">
import { onMounted } from "vue";
import { initFromCSS } from "./sync-css";
import "./sync-language";
import { language } from "@/state/language";

onMounted(() => {
  initFromCSS();
});
</script>

<template>
  <slot />
</template>
