<!-- music app -->
<template>
  <div class="music">
    <div class="music-header">
      <div class="music-header-left">
        <div class="music-header-left-icon">
          <el-button icon="Upload" circle @click.stop="uploadMusic" />
          <input
            @change="uploadFile"
            accept="audio/*"
            type="file"
            ref="uploadInput"
            style="display: none"
          />
        </div>
      </div>
      <div class="music-header-right"></div>
    </div>
    <div class="music-body">
      <div class="music-body-right">
        <div class="music-list">
          <div
            class="music-list-item"
            v-for="item in musicList"
            :key="item.path"
            @click="playMusic(item)"
          >
            <span>{{ item.name }}</span>
          </div>
        </div>

        <div class="viewer">
          <div class="music-title">
            {{ chosenMusic.path ? basename(chosenMusic.path) : "" }}
          </div>
          <div class="viewer-img">
            <div class="ani-text">
              {{ chosenMusic.path ? basename(chosenMusic.path) : "" }}
            </div>
          </div>
          <audio-player
            ref="audioplayer"
            theme-color="#444"
            :audio-list="audioList.map((elm: any) => elm.url)"
          />
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { OsFileWithoutContent, join, useSystem, basename } from "@/system";
import AudioPlayer from "@liripeng/vue-audio-player";
//import FileUploader from '@/components/apps/FileUploader.vue';

const sys = useSystem();
const musicList = ref<OsFileWithoutContent[]>([]);
const uploadInput = ref();
onMounted(async () => {
  refershFileLst();
});
async function refershFileLst() {
  const list = await sys.fs.readdir(join(sys._options.userLocation || "", "Music"));
  musicList.value = list.filter((item: any) => {
    return item.path.endsWith(".mp3");
  });
}

function uploadMusic() {
  uploadInput.value.click();
}
function uploadFile(ev: Event) {
  const tar = ev.target as HTMLInputElement;

  if (tar.files) {
    const reader: any = new FileReader();
    reader.readAsArrayBuffer(tar.files[0]);
    reader.onloadend = function () {
      if (tar.files) {
        const file = tar.files[0];
        sys.fs
          .writeFile(
            join(sys._options.userLocation || "", "Music", file.name),
            reader.result
          )
          .then(() => {
            tar.value = "";
            refershFileLst();
            sys.createNotify({
              title: "上传成功",
              content: "上传成功",
            });
          });
      }
    };
  }
}

function base64ToBlobUrl(base64: string) {
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
const audioplayer = ref<any>(null);
const audioList = ref<any>([]);
const chosenMusic = ref<any>({});
async function playMusic(item: OsFileWithoutContent) {
  chosenMusic.value = item;
  const fileC: any = await sys.fs.readFile(item.path);
  let content;
  if (fileC) {
    if (typeof fileC === "string") {
      content = base64ToBlobUrl(fileC.replace(/^data:(.)*;base64,/, ""));
    }
    if (fileC instanceof ArrayBuffer) {
      // 创建一个Blob对象，传入ArrayBuffer和对应的MIME类型
      const blob = new Blob([fileC], { type: "application/pdf" });
      // 使用URL.createObjectURL方法创建Blob URL
      content = URL.createObjectURL(blob);
    }

    audioList.value = [
      {
        name: "audio 1",
        url: content,
      },
    ];
    setTimeout(() => {
      audioplayer?.value?.play();
    }, 200);
  } else {
    audioList.value = [];
    chosenMusic.value.path = "无法打开文件";
  }
}
</script>

<style scoped lang="scss">
$primary-color: #0078d4; // Win11主题色
$neutral-light: #f3f2f1; // 浅中性色
$neutral-dark: #1e1e1e; // 深中性色
$border-radius: 8px; // Win11常用圆角大小
$shadow-depth-1: 0 1px 2px 0 rgba(0, 0, 0, 0.2), 0 1px 4px 0 rgba(0, 0, 0, 0.16); // 浅阴影
$shadow-depth-2: 0 2px 4px 0 rgba(0, 0, 0, 0.24), 0 1px 8px 0 rgba(0, 0, 0, 0.16); // 深阴影

.music {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  user-select: none;
  background-color: $neutral-light;
  border-radius: $border-radius;
  box-shadow: $shadow-depth-1;
  .music-header {
    width: 100%;
    // height: 50px;
    display: flex;
    flex: 0 0 50px; /* 固定高度为100px */
    flex-shrink: 0;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
    border-top-left-radius: inherit;
    border-top-right-radius: inherit;
    .music-header-left {
      width: 50%;
      height: 100%;
      display: flex;
      justify-content: flex-start;
      align-items: center;
      .music-header-left-icon {
        padding-left: 10px;
        display: flex;
        justify-content: center;
        align-items: center;
      }
      .music-header-left-title {
        width: 100px;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        span {
          font-size: 16px;
          color: #000;
        }
      }
    }
    .music-header-right {
      width: 50%;
      height: 100%;
      display: flex;
      justify-content: flex-end;
      align-items: center;
      .music-header-right-icon {
        width: 30px;
        height: 30px;
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }
  }
  .music-body {
    width: 100%;
    height: 0;
    flex: 1;
    display: flex;
    justify-content: flex-start;
    align-items: center;

    .music-body-right {
      width: calc(100% - 80px);
      height: 100%;
      display: flex;
      gap: 16px;
      .music-list {
        display: flex;
        flex-direction: column;
        height: 100%;
        width: 30%;
        padding-top: 20px;
        overflow-y: auto;
        border-right: 1px solid rgba(0, 0, 0, 0.1);
        .music-list-item {
          padding-left: 20px;
          cursor: pointer;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          flex-shrink: 0;
          margin: 2px;
          padding: 4px;
          border-radius: 4px;
          transition: all 0.1s;
        }
        .music-list-item:hover {
          color: white;
          background-color: #808080;
        }
      }
    }
  }
}

.viewer {
  width: 300px;
  height: 100%;
  position: relative;
  top: 20%;
  margin: 0 auto;
}
.viewer-img {
  /* 文本圆形排列，并一直旋转 */
  position: relative;
  width: 60px;
  height: 60px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 50%;
  border: 1px solid #ccc;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
  animation: rotate 10s linear infinite;
  /* 文字滚动 */
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  box-shadow: 0 0 30px 2px #c8c8c8, 0 0 1px 25px #c0c0c0 inset, 0 0 5px 35px #545454 inset,
    0 0 1px 40px #000000 inset;
  transition: all 0.3s;
}
.viewer-img:hover {
  box-shadow: 0 0 30px 4px #c8c8c8, 0 0 1px 25px #c0c0c0e4 inset,
    0 0 5px 35px #6f6f6fd8 inset, 0 0 1px 40px rgb(0, 0, 0) inset;
}
.viewer-img::after {
  content: " ";
  display: block;
  position: absolute;
  width: 20px;
  height: 20px;
  background-color: #ffffff;
  border-radius: 50%;
  margin: 0 auto;
}
@keyframes rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
.ani-text {
  display: inline-block;
  white-space: nowrap;
  animation: 4s wordsLoop linear infinite normal;
}

@keyframes wordsLoop {
  0% {
    transform: translateX(100%);
  }
  100% {
    transform: translateX(-100%);
  }
}
</style>
