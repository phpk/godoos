<template>
  <div class="viewer">
    <embed
      id="mainpdf"
      :src="content"
      type="application/pdf"
      width="100%"
      height="100%"
    />
  </div>
</template>
<script setup lang="ts">
import { BrowserWindow } from "@/system";
import { inject } from "vue";
const window: BrowserWindow | undefined = inject("browserWindow");
const val = window?.config.content;
let content: any;
if (typeof val === "string") {
  content = base64PDFToBlobUrl(val.replace(/^data:(.)*;base64,/, ""));
}
if (val instanceof ArrayBuffer) {
  // 创建一个Blob对象，传入ArrayBuffer和对应的MIME类型
  const blob = new Blob([val], { type: "application/pdf" });
  // 使用URL.createObjectURL方法创建Blob URL
  content = URL.createObjectURL(blob);
}
//const content = base64PDFToBlobUrl(window?.config.content.replace(/^data:(.)*;base64,/, ''));
function base64PDFToBlobUrl(base64: string) {
  const binStr = atob(base64);
  const len = binStr.length;
  const arr = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    arr[i] = binStr.charCodeAt(i);
  }
  const blob = new Blob([arr], { type: "application/pdf" });
  const url = URL.createObjectURL(blob);
  return url;
}
</script>
<style scoped>
.viewer {
  width: 100%;
  height: 100%;
  background-color: #ffffff;
}
.viewer-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
