<script setup lang="ts">
import { System } from "@/system/index.ts";
import { onMounted, computed } from "vue";
// 使用ref和computed改进状态管理
const isInstall: any = computed(() => {
  const item = localStorage.getItem("isGodoOSInstall");
  // 显式地进行类型检查和转换
  return item === "true";
});

onMounted(() => {
  new System();
});
// router.beforeEach((to, _, next) => {
//   console.log(to.name);
//   if (!isInstall.value && to.name !== "install") {
//     next("/install");
//   } else {
//     next();
//   }
// });
</script>

<template>
  <screen />
</template>
