import { ref, watch } from "vue";
import { i18n } from "@/runtime/i18n/i18n";

export const language = ref<"en" | "ru">("ru");

watch(
  language,
  (newLang) => {
    i18n.global.locale.value = newLang;
  },
  { immediate: true },
);
