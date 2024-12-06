<template>
  <iframe class="setiframe" allow="fullscreen" ref="storeRef" :src="src"></iframe>
</template>
<script lang="ts" setup name="IframeFile">
//@ts-ignore
import { BrowserWindow, Dialog, Notify, System } from "@/system";
import { notifyError } from "@/util/msg";
import { getSplit } from "@/system/config";
import { base64ToBuffer, isBase64 } from "@/util/file";
import { generateRandomString } from "@/util/common";
import { inject, onMounted, onUnmounted, ref, toRaw } from "vue";
import { askAi } from "@/hook/useAi";
import { useChooseStore } from "@/stores/choose";
import eventBus from "@/system/event/eventBus";
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
const storeRef = ref<HTMLIFrameElement | null>(null);
let hasInit = false;

const choose = useChooseStore();
const componentID = generateRandomString(16);
// console.log("唯一ID：", componentID);
const saveFile = async (e: any) => {
  if (e.componentID !== componentID) return;
  let data = JSON.parse(e.eventData.data);
  let ext: any = props.ext;
  if (ext instanceof Array) {
    ext = ext[0];
  }
  if (data.ext) {
    ext = data.ext;
  }
  const fileName = e.fileName == '' ? data.title : e.fileName
  let path: string
  e.filePath !== ""
    ? (path = `${e.filePath}/${fileName}.${ext}`)
    : (path = `${SP}C${SP}Users${SP}Desktop${SP}${fileName}.${ext}`);
  // console.log('路径：',path);
  await writeFile(path, data, e.fileName, true);
};
eventBus.on("saveFile", saveFile);
const writeFile = async (path: string, data: any, title: string, isNewFile: boolean) => {
  if (await sys?.fs.exists(path) && isNewFile) {
    let res = await Dialog.showMessageBox({
      type: "info",
      title: "提示",
      message: "存在相同的文件名-" + title,
      buttons: ["覆盖文件", "取消"],
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
  let res = await sys?.fs.writeFile(path, data.content);
  if (res.code === -1 && res.error == "needPwd") {
    const temp = await Dialog.showInputBox()
    if (temp.response !== 1) {
      return
    }
    const header:any = {}
    header.pwd = temp?.inputPwd ? temp?.inputPwd : ''
    res = await sys?.fs.writeFile(path, data.content, header);
    if (res.code === -1) {
      notifyError(res.message)
      return
    }
  }
  // console.log("编写文件：", res, isShare);
  new Notify({
    title: "提示",
    content: res.message,
    // content: res.code === 0 ? "文件已保存" : res.message,
  });
  sys.refershAppList();
};
const eventHandler = async (e: MessageEvent) => {
  const eventData = e.data;
  // console.log('是否同一个：', componentID == eventData.componentID);
  if (eventData.type == props.eventType) {
    // if (eventData.componentID !== componentID) return;
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
    if (title.indexOf("." + ext) > -1) {
      title = title.replace("." + ext, "");
    }

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
      choose.saveFile("选择地址", "*", componentID, eventData, ext);
      return;
      // path = `${SP}C${SP}Users${SP}Desktop${SP}${title}.${ext}`;
    }
    writeFile(path, data, title, false);
  } else if (eventData.type == "initSuccess") {
    if (hasInit) {
      return;
    }
    hasInit = true;
    let content = win?.config?.content;
    let title = win.getTitle();
    //console.log("win.config;", win?.config);
    // console.log(title);
    title = title.split(SP).pop();
    content = toRaw(content);
    if(typeof content == "string"){
      content = content.trim()
    }

    if (content && content !== "") {
      storeRef.value?.contentWindow?.postMessage(
        {
          type: "init",
          data: { content, title },
          componentID,
        },
        "*"
      );
    } else {
      storeRef.value?.contentWindow?.postMessage(
        {
          type: "start",
          title,
          componentID,
        },
        "*"
      );
    }
  } else if (eventData.type == "close") {
    // console.log("关闭");
    win.close();
  } else if (eventData.type == "saveMind") {
    // console.log("保存");
    const data = eventData.data;
    const path = win?.config?.path;
    //console.log(path,data)
    const winMind = new BrowserWindow({
      title: data.title,
      url: "/mind/index.html",
      frame: true,
      config: {
        ext: "mind",
        path: path,
        content: data.content,
      },
      icon: "gallery",
      width: 700,
      height: 500,
      x: 100,
      y: 100,
      //center: true,
      minimizable: false,
      resizable: true,
    });
    winMind.show();
  } else if (eventData.type == "aiCreater") {
    console.log("传递内容： ", eventData);
    let postData: any = {};
    if (eventData.data) {
      postData.content = eventData.data;
    }
    if (eventData.title) {
      postData.title = eventData.title;
    }
    if (eventData.category) {
      postData.category = eventData.category;
    }
    //console.log(postData,eventData.action)
    // 模拟AI返回数据
    const res: any = await askAi(postData, eventData.action);
    storeRef.value?.contentWindow?.postMessage(
      {
        type: "aiReciver",
        data: res,
        action: eventData.action,
      },
      "*"
    );
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
