<template>
  <div
    class="window-outer"
    :class="{
      focus: focusState && currentRouter !== 'main',
    }"
  >
    <div class="upbar" v-dragable>
      <div class="upbar-left">
        <div class="back-arr" v-if="currentRouter !== 'main'" @click="back">←</div>
        <div class="upbar-text">
          {{ t("setting") }}
        </div>
      </div>
      <div class="upbar-right">
        <WinUpButtonGroup :browser-window="browserWindow"></WinUpButtonGroup>
      </div>
    </div>
    <Transition v-for="item in setList" :key="item.key" name="fade" appear>
      <component :is="stepComponent(item.content)" v-if="currentRouter === item.key" />
    </Transition>
    <Transition name="fade" appear>
      <div class="outer" v-if="currentRouter === 'main'">
        <div class="uper_tab">
          <div class="tab">
            {{ t("windows.setting") }}
          </div>
        </div>
        <div class="outer_main">
          <div class="main_uper">
            <div
              class="set_item"
              v-for="item in setList"
              :key="item.title"
              @click="openSet(item.key)"
              v-glowing="{
                color: '#3c3c3ce4',
                scale: 0.6,
              }"
            >
              <div class="set_item-img">
                <!-- <img class="set_item-img-img" :src="item.icon" /> -->
                <svg class="icon" aria-hidden="true" style="font-size: 2em">
                  <use :xlink:href="'#icon-' + item.icon"></use>
                </svg>
              </div>
              <div class="set_item-right">
                <div class="set_item-title">{{ item.title }}</div>
                <div class="set_item-desc">{{ item.desc }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>
<script lang="ts" setup>
import { inject, ref } from "vue";

import { BrowserWindow } from "@/system/window/BrowserWindow";
import { t } from "@/i18n";
import { useSystem } from "@/system";
import { vDragable } from "@/system/window/MakeDragable";
import { vGlowing } from "@/util/glowingBorder";
import { stepComponent } from "@/util/stepComponent";
const browserWindow = inject<BrowserWindow>("browserWindow")!;
const sys = useSystem();
const currentRouter = ref(browserWindow.config?.router || "main");

const focusState = ref(false);
browserWindow?.on("focus", () => {
  focusState.value = true;
});
browserWindow?.on("blur", () => {
  focusState.value = false;
});
function back() {
  currentRouter.value = "main";
}
function openSet(key: string) {
  currentRouter.value = key;
}
const setList = ref([
  {
    key: "system",
    title: t("system"),
    desc: "存储、备份还原、用户角色",
    icon: "system",
    content: "SetSystem",
  },

  {
    key: "custom",
    title: "代理",
    desc: '代理/远程桌面',
    icon: "personal",
    content: "SetCustom",
  },

  {
    key: "language",
    title: '语言',
    desc: t("language"),
    icon: "language",
    content: "SetLang",
  },

  {
    key: "account",
    title: "屏幕",
    desc: '壁纸/锁屏/广告',
    icon: "account",
    content: "SetAccount",
  },

  ...(sys._rootState.settings ? sys._rootState.settings : []),
]);
</script>

<style lang="scss" scoped>
@import "@/assets/main.scss";
.window-outer {
  background-color: white;
  width: 100%;
  height: 100%;
  border: #0076d795 1px solid;
  box-sizing: border-box;
  // box-shadow: inset -599px 0px 0px 0px #ffffff;
  transition: background-color 0.1s;
  overflow: hidden;
}

.window-outer.focus {
  background-color: rgba(255, 255, 255, 0.704);
  backdrop-filter: blur(10px);
}
.upbar {
  height: 40px;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  user-select: none;
  .upbar-left {
    display: flex;
    align-items: center;
    font-size: 12px;
    height: 100%;
    .back-arr {
      padding: 0 6px;
      margin-right: 3px;
      font-size: 14px;
      width: 40px;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
    }
    .back-arr:hover {
      color: var(--color-ui-gray);
      background-color: var(--color-dark);
    }
    .upbar-text {
      font-size: 14px;
      margin-left: 10px;
    }
  }
  .upbar-right {
    width: calc(100% - 200px);
    display: flex;
    justify-content: flex-end;
    height: 100%;
    background-color: white;
  }
}

.outer {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background-color: #ffffff;
  user-select: none;
  /* background-color: rgb(241, 241, 241); */
}

.outer_main {
  height: 100%;
  margin: 0px 10px 00px 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.main_uper {
  display: flex;
  flex-direction: row;
  justify-content: center;
  flex-wrap: wrap;
  align-items: center;
  height: 50px;
  width: 90%;
}

.set_item {
  margin: 10px 20px;
  padding: 2px 8px 10px 2px;
  width: 200px;
  display: flex;
  align-items: center;
  background-color: white;
  border: 1px solid #99999900;
  transition: all 0.1s;
  position: relative;
  z-index: 2;

  .set_item-img {
    width: 30px;
    margin: 10px 16px;
    flex-shrink: 0;
    img {
      width: 100%;
    }
  }
  .set_item-right {
    flex: 1;
  }

  .set_item-title {
    padding: 4px 0px;
  }

  .set_item-desc {
    font-size: 10px;
    color: #999999;
  }
}
.set_item:hover {
  border: 1px solid #99999954;
  box-shadow: 0 0 0 1px #999999 inset;
}
.uper_tab {
  /* width: 90%; */
  margin: 10px 10px -1px 10px;
  font-size: 12px;

  z-index: 1;
  display: flex;
}

.tab {
  text-align: center;
  width: 100%;
  height: 40px;
  font-size: 20px;
  font-weight: 200;
  /* border: 1px solid #d9d9d9; */
  border-bottom: none;
}

.tab_unactive {
  text-align: center;
  width: 40px;
  background-color: rgb(241, 241, 241);

  border: 1px solid #d9d9d9;
  /* border-bottom:none; */
}

.bottom {
  flex-shrink: 0;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  flex-shrink: 0;
  width: 90%;
  margin: 6px auto;
  // background-color: #f0f0f0;
}

.bottom_button {
  width: 80px;
  height: 26px;
  line-height: 26px;
  text-align: center;
  font-size: 13px;
  background-color: #e1e1e1;
  border: 1px solid #999999;
  transition: all 0.2s;
  box-sizing: border-box;
  margin: 0 0 0 10px;
}

.bottom_button:first-child {
  border: 1px solid #0078d7;
  box-shadow: 0 0 0px 1px #0078d7 inset;
}

.bottom_button:hover {
  border: 1px solid #0078d7;
  background-color: #e5f1fb;
  box-shadow: 0 0 0 0px #0078d7 inset;
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s ease;
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
