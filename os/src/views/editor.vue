<template>
    <IframeApp :win="props.win" />
</template>
<script setup lang="ts">
import { defineProps } from "vue";
const props = defineProps({
  win: {
    type: Object,
    required: true,
  },
});
</script>