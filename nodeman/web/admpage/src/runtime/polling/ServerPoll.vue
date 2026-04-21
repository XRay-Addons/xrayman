<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import { startServerPoll, stopServerPoll } from "./server-poll";

onMounted(() => startServerPoll());
onUnmounted(() => stopServerPoll());
</script>

<template></template>
