<template>
  <a-config-provider :theme="theme">
    <slot />
  </a-config-provider>
</template>

<script setup>
const theme = {
  token: {
    colorPrimary: "#6366f1",
    colorBgBase: "#020617",
    colorText: "#e5e7eb",
  },
};
</script>
