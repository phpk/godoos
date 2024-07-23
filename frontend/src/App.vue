<script setup lang="ts">
import { System } from "@/system/index.ts";
import { onMounted } from "vue";

onMounted(() => {
  new System();
});
</script>

<template>
  <screen />
</template>
