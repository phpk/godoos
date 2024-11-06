<template>
  <div class="outer">
    <el-input class="win-input" v-model="name" />
    <WinButton @click="confirm">{{ t("confirm") }}</WinButton>
  </div>
</template>
<script setup lang="ts">
import { inject, ref } from "vue";
import { emitEvent } from "@/system/event";
import {
  Dialog,
  t,
  Notify,
  basename,
  join,
  BrowserWindow,
  useSystem,
  OsFileWithoutContent,
} from "@/system";
import { notifyError } from "@/util/msg";

const browserWindow: BrowserWindow = inject("browserWindow")!;
const name = ref(basename((browserWindow.config.content as OsFileWithoutContent).path));

function confirm() {
  if (name.value.length > 40) {
    new Notify({
      title: "提示",
      content: "文件名过长",
    });
    return;
  }
  let oldPath = ''
  if(browserWindow.config.content?.isShare) {
    oldPath = browserWindow.config.content.parentPath
  } else {
    const temp = browserWindow.config.content.path.split('/')
    temp.pop()
    oldPath = temp.join('/')
  } 
  const newPath = join(oldPath, name.value);
  useSystem()
    ?.fs.rename(browserWindow.config.content.path, newPath)
    .then((res: any) => {
      if(!res || res.code === -1) {
        notifyError(res.message || '改名失败')
        return 
      }
      emitEvent("file.props.edit");
      browserWindow.emit("file.props.edit", newPath);
      browserWindow.close();
    })
    .catch((e: any) => {
      Dialog.showMessageBox({
        message: e,
        type: "error",
      });
    });
}
</script>
<style lang="scss" scoped>
.outer {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;

  .win-input {
    font-size: 20px;
    width: 200px;
    height: 40px;
    margin-bottom: 40px;
    outline: none;
    border: 1px solid black;

    &:focus {
      border: 1px solid var(--color-blue);
    }
  }
}
</style>
@/system/core/FileSystem
