<template>
  <div class="outer">
    <div class="uper">
      <div class="up-text">应用商店</div>
    </div>
    <div class="store-handle" v-dragable>
      <div v-if="!closing" @click="closeWin" class="close-button">×</div>
    </div>
    <div class="main">
      <div class="left">
        <div class="left-icon">
          <div class="icon-derc"></div>
          <svg
            t="1694613650127"
            class="icon"
            viewBox="0 0 1024 1024"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
            p-id="10617"
            width="200"
            height="200"
          >
            <path
              d="M511.813 92.188L92.188 428.844v502.5h295.406V647.469h248.344v283.875h295.406v-502.5z"
              p-id="10618"
            ></path>
          </svg>
        </div>
      </div>
      <div class="store">
        <div v-if="isready" class="store-top">
          <!-- <div class="left-bar"></div> -->
          <div class="right-main">
            <div class="main-title">
              <!-- <span class="sub-title">热门应用 </span> -->
            </div>
            <div class="swiper">
              <div class="swiper-txt">主页</div>
              <div class="swiper-inner">
                <div class="swiper-tab">
                  <img src="/image/store/banner1.jpg" />
                </div>
                <div class="swiper-tab">
                  <img src="/image/store/banner2.jpg" />
                </div>
                <div class="swiper-tab">
                  <img src="/image/store/banner3.jpg" />
                </div>
                <div class="swiper-tab">
                  <img src="/image/store/banner2.jpg" />
                </div>
              </div>
            </div>
            <div class="main-app">
              <div v-for="item in storeList" class="store-item" :key="item.name">
                <AppItem
                  :item="item"
                  :installed-list="installedList"
                  :install="install"
                  :uninstall="uninstall"
                />
              </div>
            </div>
          </div>
        </div>
        <div v-else class="store-noready">
          <div id="wait">
            <div class="waitd" id="wait1"></div>
            <div class="waitd" id="wait2"></div>
            <div class="waitd" id="wait3"></div>
            <div class="waitd" id="wait4"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { inject, onMounted, ref } from "vue";
import { BrowserWindow, System, Dialog, vDragable } from "@/system";
import { i18n } from "@/i18n";
import storeList from "@/assets/store.json";
const browserWindow: BrowserWindow = inject("browserWindow")!;
import { getSystemKey, setSystemKey } from "@/system/config";
const sys: any = inject<System>("system");
console.log(storeList);
const isready = ref(false);
const installed = getSystemKey("intstalledPlugins");
const installedList: any = ref(installed);
const closing = ref(false);
function closeWin() {
  closing.value = true;
  setTimeout(() => {
    browserWindow.close();
  }, 200);
}
onMounted(() => {
  setTimeout(() => {
    if (!isready.value) {
      isready.value = true;
    }
  }, 1000);
});
function setCache() {
  setSystemKey("intstalledPlugins", installedList.value);
  //localStorage.setItem("godoOS_installedApp", JSON.stringify(installedList.value))
  setTimeout(() => {
    sys.refershAppList();
  }, 1000);
}
function install(item: any) {
  //console.log(item)
  sys.fs.writeFile(
    `${sys._options.userLocation}Desktop/${item.name}.url`,
    `link::url::${item.url}::${item.icon}`
  );
  Dialog.showMessageBox({
    message: i18n("install.success"),
    type: "info",
    buttons: [i18n("confirm")],
  });
  installedList.value.push(item.name);
  setCache();
}

function uninstall(item: any) {
  sys.fs.unlink(`${sys._options.userLocation}Desktop/${item.name}.url`);
  Dialog.showMessageBox({
    message: i18n("uninstall.success"),
    type: "info",
    buttons: [i18n("confirm")],
  });
  delete installedList.value[installedList.value.indexOf(item.name)];
  setCache();
}
</script>

<style>
/*定义滚动条高宽及背景
 高宽分别对应横竖滚动条的尺寸*/
::-webkit-scrollbar {
  width: 5px;
  height: 16px;
  background-color: #ffffff;
}
/*定义滚动条轨道
  内阴影+圆角*/
::-webkit-scrollbar-track {
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  border-radius: 10px;
  background-color: #ffffff;
}
/*定义滑块
  内阴影+圆角*/
::-webkit-scrollbar-thumb {
  border-radius: 10px;
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
  background-color: rgba(132, 132, 132, 0.537);
}
</style>
<style scoped>
.outer {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.uper {
  padding: 0px 20px;
  font-size: 12px;
  height: 40px;
  line-height: 40px;
  flex-shrink: 0;
  color: rgb(134, 134, 134);
  background-color: rgb(243, 243, 243);
}
.store-handle {
  width: 100%;
  height: 40px;
  user-select: none;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  position: absolute;
  top: 10px;
  right: 20px;
}

.main {
  width: 100%;
  display: flex;
  flex: 1;
  overflow: hidden;
  background-color: rgb(243, 243, 243);
}
.left {
  width: 60px;
  height: 60px;
  flex-shrink: 0;
  background-color: rgb(243, 243, 243);
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
}
.left-icon {
  width: 56px;
  height: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: white;
  border-radius: 6px;
  position: relative;
}
.icon-derc {
  height: 20px;
  width: 4px;
  border-radius: 4px;
  background-color: #363533;
  position: absolute;
  left: 2px;
}
.left-icon svg {
  width: 20px;
  height: 20px;
  color: #363533;
  stroke: #363533;
}
.store {
  /* width: 100%; */
  flex: 1;
  background-color: rgb(255, 255, 255);
  border-top-left-radius: 10px;
  overflow: hidden;
}

.store-top {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  overflow: hidden;
}
.left-bar {
  width: 60px;
  flex-shrink: 0;
}
.right-main {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
}

.main-title {
  font-size: 26px;
  font-weight: bold;
  font-variant: small-caps;
  font-variant-ligatures: discretionary-ligatures;
  margin: 8px;
  margin-top: 8px;
  user-select: none;
  overflow: hidden;
  white-space: nowrap;
}
.sub-title {
  padding-top: 30px;
  font-size: 20px;
  font-weight: bold;
  margin: 10px;
  user-select: none;
  overflow: hidden;
  white-space: nowrap;
}
.sub-tip {
  font-size: 12px;
  color: rgb(11, 31, 111);
  margin-left: 10px;
  user-select: none;
}
.swiper {
  width: max-content;
  height: 300px;
  /* overflow: hidden; */
  margin: 10px;
  margin-top: 0px;
  position: relative;
}
@keyframes swiperAni {
  0%,
  33% {
    transform: translateX(0px);
  }

  36%,
  66% {
    transform: translateX(-600px);
  }

  69%,
  96% {
    transform: translateX(-1200px);
  }
  100% {
    transform: translateX(0px);
  }
}
.swiper-inner {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 10px;
  animation: swiperAni 20s ease-in-out infinite;
}
.swiper-tab {
  width: 600px;
  height: 300px;
  background-color: rgb(243, 243, 243);
  border-radius: 20px;
  box-shadow: 0px 10px 20px 1px #2524241f;
}
.swiper-tab img {
  width: 100%;
  height: 100%;
  border-radius: 20px;
  object-fit: cover;
}
.swiper-txt {
  position: absolute;
  top: 20px;
  left: 30px;
  font-weight: 600;
  color: white;
  z-index: 10;
  text-shadow: 0px 0px 5px #00000058;
}
.main-app {
  width: 100%;
  /* height: 100%; */
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: flex-start;
  align-items: flex-start;
  padding: 10px;
  overflow: auto;
  align-content: flex-start;
}
.store-noready {
  width: 100%;
  height: 100%;
  background-color: #cdcdcd69;
}

#wait {
  position: absolute;
  left: 50%;
  top: calc(50% + 150px);
}

.waitd {
  position: absolute;
  width: 10px;
  height: 10px;
  left: 30px;
  background-color: azure;
  border-radius: 50%;
  transform-origin: -15px 0;
}

#wait1 {
  animation: dotAni1 2s linear infinite;
}

#wait2 {
  animation: dotAni2 2s linear infinite;
}

#wait3 {
  animation: dotAni3 2s linear infinite;
}

#wait4 {
  animation: dotAni4 2s linear infinite;
}

@keyframes dotAni1 {
  0% {
    transform: rotateZ(0deg);
  }

  20% {
    transform: rotateZ(240deg);
  }

  85% {
    transform: rotateZ(290deg);
  }

  100% {
    transform: rotateZ(360deg);
  }
}

@keyframes dotAni2 {
  0% {
    transform: rotateZ(0deg);
  }

  35% {
    transform: rotateZ(240deg);
  }

  85% {
    transform: rotateZ(290deg);
  }

  100% {
    transform: rotateZ(360deg);
  }
}

@keyframes dotAni3 {
  0% {
    transform: rotateZ(0deg);
  }

  50% {
    transform: rotateZ(240deg);
  }

  85% {
    transform: rotateZ(290deg);
  }

  100% {
    transform: rotateZ(360deg);
  }
}

@keyframes dotAni4 {
  0% {
    transform: rotateZ(0deg);
  }

  65% {
    transform: rotateZ(240deg);
  }

  85% {
    transform: rotateZ(290deg);
  }

  100% {
    transform: rotateZ(360deg);
  }
}
</style>
