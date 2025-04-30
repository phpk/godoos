<template>
  <div
    class="rect"
    draggable="false"
    ref="parentRef"
    :style="{
      left: rect.left + 'px',
      top: rect.top + 'px',
      width: rect.width + 'px',
      height: rect.height + 'px',
    }"
  ></div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';

defineProps<{
  rect: {
    left: number;
    top: number;
    width: number;
    height: number;
  };
}>();

const parentRef = ref<HTMLElement>();
const parentRect = ref<DOMRect | null>(null);

onMounted(() => {
  const parent = parentRef.value?.parentElement;
  if (parent) {
    parentRect.value = parent.getBoundingClientRect();
  }
});
</script>

<style lang="scss" scoped>
.rect {
  position: absolute;
  border: 1px dashed #007bff;
  background: rgba(0, 123, 255, 0.1);
  pointer-events: none; // 确保框选区域不影响其他元素的交互
  z-index: 9999;
}
</style>