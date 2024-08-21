<template>
  <iframe class="setiframe" allow="fullscreen" ref="storeRef" :src="src"></iframe>
</template>
<script lang="ts" setup name="IframeFile">
//@ts-ignore
import { BrowserWindow, Notify, System, Dialog } from "@/system";
import { ref, onMounted, inject, onUnmounted, toRaw } from "vue";
import { isBase64, base64ToBuffer } from "@/util/file";
import { getSplit } from "@/system/config";
const SP = getSplit();

const sys: any = inject<System>("system");
const win: any = inject<BrowserWindow>("browserWindow");
const props = defineProps({
  src: {
    type: String,
    default: "",
  },
  eventType: {
    type: String,
    default: "",
  },
  ext: {
    type: String,
    default: "md",
  },
});
//console.log(props);
//let path = win?.config?.path;
const storeRef = ref<HTMLIFrameElement | null>(null);
let hasInit = false;
const eventHandler = async (e: MessageEvent) => {
  const eventData = e.data;
  // console.log(path)
  //console.log(eventData);
  if (eventData.type == props.eventType) {
    let data = JSON.parse(eventData.data);
    let title = data.title;
    let path;
    let ext: any = props.ext;
    if (ext instanceof Array) {
      ext = ext[0];
    }
    if (data.ext) {
      ext = data.ext;
    }
    // console.log(ext)
    // console.log(data)
    if (win.config && win.config.path) {
      path = win.config.path;
      let fileTitleArr = path.split(SP).pop().split(".");
      let oldExt = fileTitleArr.pop();
      let fileTitle = fileTitleArr.join(".");
      if (fileTitle != title) {
        path = path.replace(fileTitle, title);
      }
      if (oldExt != ext) {
        path = path.replace("." + oldExt, "." + ext);
      }
    } else {
      path = `${SP}C${SP}Users${SP}Desktop${SP}${title}.${ext}`;
    }

    if (await sys?.fs.exists(path)) {
      let res = await Dialog.showMessageBox({
        type: "info",
        title: "提示",
        message: "存在相同的文件名-" + title,
        buttons: ["覆盖文件?", "取消"],
      });
      //console.log(res)
      if (res.response > 0) {
        return;
      }
    }
    if (typeof data.content === "string") {
      if (data.content.indexOf(";base64,") > -1) {
        const parts = data.content.split(";base64,");
        data.content = parts[1];
      }
      if (isBase64(data.content)) {
        data.content = base64ToBuffer(data.content);
        //console.log(data.content)
      }
    }
    //console.log(data.content)
    await sys?.fs.writeFile(path, data.content);
    new Notify({
      title: "提示",
      content: "文件已保存",
    });
    sys.refershAppList();
  } else if (eventData.type == "initSuccess") {
    if (hasInit) {
      return;
    }
    hasInit = true;
    let content = win?.config?.content;
    console.log(win?.config);
    let title = win.getTitle();
    // console.log(title);
    title = title.split(SP).pop();
    //console.log(title);
    if (!content && win?.config.path) {
      content = await sys?.fs.readFile(win?.config.path);
    }
    content = toRaw(content);
    // console.log(content);
    if (content && content !== "") {
      storeRef.value?.contentWindow?.postMessage(
        {
          type: "init",
          data: { content, title },
        },
        "*"
      );
    } else {
      storeRef.value?.contentWindow?.postMessage(
        {
          type: "start",
          title,
        },
        "*"
      );
    }
  }
};
onMounted(() => {
  window.addEventListener("message", eventHandler);
});

onUnmounted(() => {
  window.removeEventListener("message", eventHandler);
});
</script>
<style scoped>
.setiframe {
  width: 100%;
  height: 100%;
  border: none;
}
</style>
