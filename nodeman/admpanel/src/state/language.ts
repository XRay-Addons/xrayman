import { ref, watch } from "vue";
import { i18n } from "@/runtime/i18n/i18n";

export const language = ref<"en" | "ru">("en");

watch(
  language,
  (newLang) => {
    i18n.global.locale.value = newLang;
  },
  { immediate: true },
);

setTimeout(() => {
  if (typeof window !== "undefined") {
    console.log("teere");
    (window as any).changeLang = (lang: "en" | "ru") => {
      language.value = lang;
    };

    (window as any).showLang = () => {
      console.log(`Current language: ${language.value}`);
      return language.value;
    };
  }
}, 2000);
