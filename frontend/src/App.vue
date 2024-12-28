<script setup lang="ts">
import { System } from "@/system/index.ts";
import { onMounted } from "vue";
// import { isMobileDevice } from "@/util/device";

onMounted(() => {
  new System();
  // isMobile.value = isMobileDevice();
  //console.log(isMobile.value);
});
// const isMobile = ref<boolean>(false);

</script>

<template>
  <screen />
  <!-- <mobile v-else></Mobile> -->
</template>
