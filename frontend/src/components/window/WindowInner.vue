<template>
  <template v-if="win.url">
    <IframeFile v-if="config.editorType == 'local'" :src="win.url" :ext="win.ext" :eventType="win.eventType" />
    <template v-if="officeFile">
      <OnlyOffice v-if="config.editorType == 'onlyoffice'" :onlyType="onlyType" :eventType="win.eventType" :ext="ext" />
    </template>
    <IframeFile v-else :src="win.url" :ext="win.ext" :eventType="win.eventType" />
  </template>

  <Suspense v-else>
    <component v-if="win.content" :is="stepComponent(win.content)" :translateSavePath="translateSavePath"
      :componentID="win.windowInfo.componentID"></component>
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
let onlyType = ref("word")
const win = ref(props.win)
//console.log(win)
const officeFile = ref(false)
if (win.value.config.path) {
  ext = win.value.config.path.split('.').pop()
  officeFile.value = isOffice(ext)
}
const translateSavePath = inject('translateSavePath')
if (props.win.path) {
  router.push(props.win.path);
}
function isOffice(ext: string) {
  const word = ['doc', 'docm', 'docx', 'docxf', 'dot', 'dotm', 'dotx', 'epub', 'fodt', 'fb2', 'mht', 'odt', 'oform', 'ott', 'oxps', 'pdf', 'rtf', 'txt', 'djvu', 'xml', 'xps'];
  const cell = ['csv', 'fods', 'ods', 'ots', 'xls', 'xlsb', 'xlsm', 'xlsx', 'xlt', 'xltm', 'xltx'];
  const slide = ['fodp', 'odp', 'otp', 'pot', 'potm', 'potx', 'pps', 'ppsm', 'ppsx', 'ppt', 'pptm', 'pptx']
  if (word.includes(ext)) {
    onlyType.value = 'word'
    return true
  } 
  else if(cell.includes(ext)){
    onlyType.value = 'cell'
    return true
  }
  else if(slide.includes(ext)){
    onlyType.value = 'slide'
    return true
  }
  else {
    return false
  }
}
</script>
