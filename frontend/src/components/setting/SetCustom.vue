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
        <h1 class="setting-title">{{ t("background") }}</h1>
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
          <ColorPicker v-model:modelValue="config.background.color" @update:modelValue="onColorChange"></ColorPicker>
        </div>
      </template>
      <template v-if="config.background.type === 'image'">
        <div class="setting-item">
          <ul class="image-gallery">
            <li
              v-for="(item, index) in config.background.imageList"
              :key="index"
              :class="config.background.url === item ? 'selected' : ''"
              @click="setBg(item)"
            >
              <img :src="item" />
            </li>
          </ul>
        </div>
        <div class="setting-item">
          <label> </label>
        </div>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { getSystemConfig, setSystemConfig } from "@/system/config";
import { ref } from "vue";
import { t, useSystem } from "@/system";
const sys = useSystem();
const items = [t("background")];
const desktopOptions = [
  {
    label: t("image"),
    value: "image",
  },
  {
    label: t("color"),
    value: "color",
  },
];

const activeIndex = ref(0);
const config:any = ref(getSystemConfig());
function setBg(item: any){
  config.value.background.url = item
  config.value.background.type = "image";
  setSystemConfig(config.value);
  sys.initBackground();

}
function onColorChange(color : string){
  config.value.background.color = color;
  config.value.background.type = "color";
  setSystemConfig(config.value);
  sys.initBackground();
}
const selectItem = (index: number) => {
  activeIndex.value = index;
};
</script>

<style scoped>
@import "./setStyle.css";
@import "@/assets/imglist.scss";
</style>
