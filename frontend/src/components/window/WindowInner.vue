<template>
  <IframeFile v-if="win.url" :src="win.url" :ext="win.ext" :eventType="win.eventType" />
  <Suspense v-else>
    <component v-if="win.content" :is="stepComponent(win.content)" :translateSavePath="translateSavePath"></component>
    <RouterView v-else />
  </Suspense>
  <!-- <component :is="window.content"></component> -->
</template>
<script setup lang="ts">
import { useRouter } from "vue-router";
import { stepComponent } from "@/util/stepComponent";
const router = useRouter();
const props = defineProps<{
  win: any;
}>();
const win = ref(props.win)
const translateSavePath = inject('translateSavePath')
if (props.win.path) {
  router.push(props.win.path);
}
</script>
