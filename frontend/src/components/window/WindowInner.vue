<template>
  <template v-if="win.url">
    <IframeFile v-if="config.editorType =='local'" :src="win.url" :ext="win.ext" :eventType="win.eventType" />
    <OnlyOffice v-else :src="win.url" :eventType="win.eventType" :ext="ext" />
  </template>
 
  <Suspense v-else>
    <component v-if="win.content" :is="stepComponent(win.content)" :translateSavePath="translateSavePath" :componentID="win.windowInfo.componentID"></component>
    <RouterView v-else />
  </Suspense>
  <!-- <component :is="window.content"></component> -->
</template>
<script setup lang="ts">
import { useRouter } from "vue-router";
import { stepComponent } from "@/util/stepComponent";
import { getSystemConfig } from "@/system/config";
const config = getSystemConfig();
const router = useRouter();
const props = defineProps<{
  win: any;
}>();
let ext = "txt"
const win = ref(props.win)
if(win.value.config.path){
  ext = win.value.config.path.split('.').pop()
}
const translateSavePath = inject('translateSavePath')
if (props.win.path) {
  router.push(props.win.path);
}
</script>
