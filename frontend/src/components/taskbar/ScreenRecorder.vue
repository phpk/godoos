<template>
  <!--截图组件-->
  <div class="state-item">
    <el-icon
      :size="18"
      class="cutRecord"
      @click="screenRecorder.actions.startRecording()"
      v-if="['idle', 'permission-requested', 'error'].includes(screenRecorder.status)"
      alt="录屏"
    >
      <Camera />
    </el-icon>
    <el-icon
      :size="18"
      class="cutRecord"
      v-if="['stopped'].includes(screenRecorder.status)"
      @click="screenRecorder.actions.resetRecording()"
      alt="录屏"
    >
      <CameraFilled />
    </el-icon>

    <el-icon
      :size="18"
      class="cutRecord"
      @click="screenRecorder.actions.stopRecording()"
      v-if="['recording', 'paused'].includes(screenRecorder.status)"
      alt="停止"
    >
      <SuccessFilled />
    </el-icon>
    <template v-if="['recording', 'paused'].includes(screenRecorder.status)">
      <el-icon
        :size="18"
        class="cutRecord"
        @click="screenRecorder.actions.resumeRecording()"
        v-if="screenRecorder.status === 'paused'"
        alt="重启"
      >
        <VideoCamera />
      </el-icon>
      <el-icon
        :size="18"
        class="cutRecord"
        @click="screenRecorder.actions.pauseRecording()"
        v-else
        alt="暂停"
      >
        <VideoPause />
      </el-icon>
    </template>
    <div class="video-container" v-if="screenRecorder.blobUrl">
      <el-button class="save-button" @click="saveRecorder">保存</el-button>
      <video
        :src="screenRecorder.blobUrl"
        class="record-video centered"
        autoplay
        controls
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSystem, Notify } from "@/system";
import useScreenRecorder from "@/util/screenRecorder";
const sys = useSystem();
const screenRecorder: any = useScreenRecorder();

async function blobToArrayBuffer(blob: Blob): Promise<ArrayBuffer> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (event) => {
      if (event.target && event.target.result instanceof ArrayBuffer) {
        resolve(event.target.result);
      } else {
        reject(new Error("Failed to read blob as ArrayBuffer."));
      }
    };
    reader.onerror = (error) => reject(error);
    reader.readAsArrayBuffer(blob);
  });
}
// 获取裁剪区域图片信息
const saveRecorder = async function () {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");
  const seconds = String(now.getSeconds()).padStart(2, "0");

  const formattedTime = year + month + day + hours + minutes + seconds;

  const title = formattedTime + "录屏";
  //console.log(title)
  const path = "/C/Users/Photo/" + title + ".mp4";
  const buffer: any = await blobToArrayBuffer(screenRecorder.blob);
  await sys.fs.writeFile(path, buffer);
  screenRecorder.actions.resetRecording();
  new Notify({
    title: "提示",
    content: "录屏已保存到图片库",
  });
};
</script>
<style>
.cutRecord {
  width: 18px;
  height: 18px;
  margin: 0px 8px;
}
.video-container {
  position: absolute;
  bottom: 60px;
  z-index: 999;
  right: 0;
  transform: translate(-50%, -50%);
  /* 保持原有的样式，如max-width, max-height等 */
  max-width: 600px;
  max-height: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.save-button {
  /* 自定义你的保存按钮样式 */
  padding: 10px;
  margin-left: auto;
}

.centered {
  /* 移除之前的固定定位，改为相对定位或静态定位，以便于居中 */
  position: static;
  max-width: 600px;
  max-height: 400px;
  /* 无需设置left和bottom */
}
</style>
