<template>
  <slot />
  <contextHolder />
</template>

<script setup lang="ts">
import { notification } from "ant-design-vue";
import { setErrorHandler } from "@/runtime/notifications/errors";

const [api, contextHolder] = notification.useNotification();

setErrorHandler((message, description) => {
  api.error({ message, description, placement: "bottomRight" });
});
</script>
