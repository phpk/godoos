<template>
    <img v-if="icon" draggable="false" :src="iconR" @error="replaceIcon" />
    <img v-else draggable="false" :src="iconimg" @error="replaceIcon" />
  </template>
  <script setup lang="ts">
  import { OsFileWithoutContent } from "@/system/core/FileSystem";
  import { useSystem } from "@/system";
  import { dealIcon } from "@/util/Icon";
  import unknownicon from '@/assets/unknown.png';
  import { ref } from "vue";
  const props = defineProps<{
    file?: OsFileWithoutContent | null;
    icon?: string;
  }>();
  
  const sys = useSystem();
  const iconimg = ref(await dealIcon(props.file, sys));
  const iconR = ref(props.icon);
  
  function replaceIcon() {
    iconR.value = unknownicon;
    iconimg.value = unknownicon;
  }
  </script>
  