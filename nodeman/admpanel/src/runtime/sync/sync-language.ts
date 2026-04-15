import { watch, onMounted } from "vue";
import { language } from "@/state/language";
import { i18n } from "@/runtime/i18n/i18n";

watch(
  language,
  (newLang) => {
    // update language in i18n
    i18n.global.locale.value = newLang;
    // pass event to static astro content
    window.dispatchEvent(new CustomEvent("language-changed"));
  },
  { immediate: true },
);
