<template>
  <a-config-provider :theme="theme">
    <a-button type="primary" @click="open = true"> Открыть форму </a-button>

    <a-modal
      v-model:open="open"
      class="modal-popup-form"
      wrap-class-name="modal-popup-form"
    >
      <slot />
    </a-modal>
  </a-config-provider>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { type Palette, PaletteItem } from "@/state/palette";
import { type ShallowRef, computed } from "vue";
import { getPaletteRef } from "@/state/palette";
import Button from "@/astro/components/primitives/_Button.astro";
import { experimental_AstroContainer } from "astro/container";

//import { getContainerRenderer as vueContainerRenderer } from "@astrojs/vue";

const open = ref(false);

function useModalTheme(paletteRef: ShallowRef<Palette>) {
  return computed(() => {
    const mainColor = paletteRef.value[PaletteItem.Title];

    return {
      components: {
        Modal: {
          colorBgElevated: mainColor,
          borderRadiusSM: 0,
          borderRadiusMD: 0,
          borderRadiusLG: 0,
          algorithm: true,
        },
      },
    };
  });
}

const theme = useModalTheme(getPaletteRef());

/*import { loadRenderers } from "astro:container";
import { experimental_AstroContainer } from "astro/container";

const renderers = await loadRenderers([vueContainerRenderer()]);
const container = await experimental_AstroContainer.create({
  renderers,
});
const result = await container.renderToString(Button);
console.log(result);*/
</script>

<style scoped>
.ant-modal-content {
  overflow: auto;
  border-radius: 24px;
}

:deep(.modal-popup-form > div > .ant-modal-content) {
  border-radius: 0px;
}

/* Стилизуем content внутри вашего модального окна */
.modal-popup-form .ant-modal-content {
  border-radius: 0;
  /* любые другие стили */
}

/* Если нужно стилизовать и другие части */
.modal-popup-form .ant-modal-header {
  border-radius: 0;
}

.modal-popup-form .ant-modal-footer {
  border-radius: 0;
}
</style>
