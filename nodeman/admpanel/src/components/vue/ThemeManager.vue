<template>
  <div class="theme-manager" style="display: contents">
    <slot />
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from "vue";
import { colors } from "./colors"; // Импортируем из отдельного файла

/*export const colors = ref({
  bgColor: "#feffa3",
  cardColor: "#d4ffea",
  titleColor: "#ffd4e5",
  buttonColor: "#eecbff",
  inputColor: "#dbdcff",
});*/

function applyTheme() {
  const root = document.documentElement;
  root.style.setProperty("--bg-color", colors.value.bgColor);
  root.style.setProperty("--card-color", colors.value.cardColor);
  root.style.setProperty("--title-color", colors.value.titleColor);
  root.style.setProperty("--button-color", colors.value.buttonColor);
  root.style.setProperty("--input-color", colors.value.inputColor);
}

watch(
  colors,
  () => {
    applyTheme();
  },
  { deep: true },
);

const changeColorTheme = (colorsArray) => {
  if (!Array.isArray(colorsArray) || colorsArray.length === 0) return;
  console.log("change thene");
  const colorKeys = Object.keys(colors.value);
  colorKeys.forEach((key, index) => {
    colors.value[key] = colorsArray[index % colorsArray.length];
  });
};

onMounted(() => {
  console.log("onMounted fired");
  window.changeColorTheme = changeColorTheme;
  applyTheme();

  // Диспатчим событие о готовности
  window.dispatchEvent(
    new CustomEvent("theme-manager-ready", {
      detail: { changeColorTheme },
    }),
  );
});

defineExpose({ changeColorTheme });
</script>
