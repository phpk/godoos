<template>
  <div class="container">
    <div class="nav">
      <ul>
        <li
          v-for="(item, index) in items"
          :key="index"
          @click="selectItem(index)"
          :class="{ active: index === activeIndex }"
        >
          {{ item }}
        </li>
      </ul>
    </div>
    <div class="setting">
      <div class="setting-item">
        <h1 class="setting-title">{{ i18n("background") }}</h1>
      </div>
      <div class="setting-item">
        <el-select v-model="config.background.type">
          <el-option
            v-for="(item, key) in desktopOptions"
            :key="key"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </div>
      <template v-if="config.background.type === 'color'">
        <div class="setting-item">
          <label> </label>
          <ColorPicker v-model="config.background.color"></ColorPicker>
        </div>
      </template>
      <template v-if="config.background.type === 'image'">
        <div class="setting-item">
          <ul class="image-gallery">
            <li
              v-for="(item, index) in config.background.imageList"
              :key="index"
              :class="config.background.url === item ? 'selected' : ''"
              @click="config.background.url = item"
            >
              <img :src="item" />
            </li>
          </ul>
        </div>
        <div class="setting-item">
          <label> </label>
        </div>
      </template>

      <div class="setting-item">
        <label> </label>
        <el-button @click="submit" type="primary">
          {{ i18n("confirm") }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { getSystemConfig, setSystemConfig } from "@/system/config";
import { ref, toRaw } from "vue";
import { Dialog, i18n, useSystem } from "@/system";
const sys = useSystem();
const items = [i18n("background")];
const desktopOptions = [
  {
    label: i18n("image"),
    value: "image",
  },
  {
    label: i18n("color"),
    value: "color",
  },
];

const activeIndex = ref(0);
const config = ref(getSystemConfig());

/** 提交背景设置 */
async function submit() {
  const val = toRaw(config.value);
  setSystemConfig(val);
  config.value = val;
  Dialog.showMessageBox({
    message: i18n("save.success"),
    title: i18n("wallpaper"),
    type: "success",
  }).then(() => {
    //location.reload();
    sys.initBackground();
  });
}

const selectItem = (index: number) => {
  activeIndex.value = index;
};
// async function submitStyle() {
//   let rootStyle = system.getConfig("rootStyle");
//   rootStyle = {
//     ...rootStyle,
//     "--icon-title-color": textColor.value,
//     "--window-border-radius": winRadius.value,
//     "--theme-main-color": taskBarColor.value,
//   };
//   await system.setConfig("rootStyle", rootStyle);

//   Dialog.showMessageBox({
//     message: i18n("save.success"),
//     title: i18n("style"),
//     type: "info",
//   });
// }
</script>

<style scoped>
@import "./setStyle.css";
@import "@/assets/imglist.scss";
</style>
