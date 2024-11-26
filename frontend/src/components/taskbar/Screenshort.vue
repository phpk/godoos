<template>
  <!--截图组件-->
  <div class="state-item">
    <!-- <svg class="icon" style="font-size:1.2em;" aria-hidden="true" @click="screenshotStatus = !screenshotStatus" alt="截屏">
      <use xlink:href="#icon-cut"></use>
    </svg> -->
    <el-icon
      class="cutImg"
      :size="18"
      @click="screenshotStatus = !screenshotStatus"
      alt="截屏"
    >
      <Scissor />
    </el-icon>
    <!-- <img :src="cutIcon" @click="screenshotStatus = !screenshotStatus" alt="截屏" class="cutImg" /> -->
    <screen-short
      v-if="screenshotStatus"
      @destroy-component="destroyComponent"
      :enableWebRtc="true"
      @get-image-data="getImg"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
//import cutIcon from '@/assets/cut.png';
const screenshotStatus = ref<boolean>(false);
import { useSystem, Notify } from "@/system";
import { isBase64, base64ToBuffer } from "@/util/file";
const sys = useSystem();
// 销毁组件函数
const destroyComponent = function (status: boolean) {
  screenshotStatus.value = status;
};
// 获取裁剪区域图片信息
const getImg = function (content: any) {
  console.log("截图组件传递的图片信息", content);

  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");
  const seconds = String(now.getSeconds()).padStart(2, "0");

  const formattedTime = year + month + day + hours + minutes + seconds;

  const title = formattedTime + "截图";
  //console.log(title)
  const path = "/C/Users/Pictures/" + title + ".png";
  // const save = {
  //   title,
  //   content,
  //   ext: 'png'
  // }
  // console.log(save)
  if (content.indexOf(";base64,") > -1) {
    const parts = content.split(";base64,");
    content = parts[1];
  }
  if (isBase64(content)) {
    content = base64ToBuffer(content);
  }
  sys.fs.writeFile(path, content);
  new Notify({
    title: "提示",
    content: "图片已保存到图片库",
  });
};
</script>
<style>
.cutImg {
  width: 18px;
  height: 18px;
  margin: 0px 8px;
}
</style>
