<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import {
  startI18nObserver,
  stopI18nObserver,
  updateDOMLanguage,
} from "./i18n-dom";

onMounted(() => {
  updateDOMLanguage();
  startI18nObserver();
});
onUnmounted(() => stopI18nObserver());
</script>

<template></template>
