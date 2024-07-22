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
      <div v-if="0 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">{{ i18n("language") }}</h1>
        </div>
        <div class="setting-item">
          <label></label>
          <el-select v-model="modelvalue">
            <el-option
              v-for="(item, key) in langList"
              :key="key"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <!-- <WinSelect
            v-model="modelvalue"
            :options="[
              {
                label: '中文',
                value: 'zh',
              },
              {
                label: 'English',
                value: 'en',
              },
            ]"
            :placeholder="i18n('please.select')"
          >
          </WinSelect> -->
        </div>

        <div class="setting-item">
          <label></label>
          <el-button @click="submit" type="primary">
            {{ i18n("confirm") }}
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { Dialog } from "@/system";
import { i18n, currentLang, setLang } from "@/i18n";

const langList = [
  {
    label: "中文",
    value: "zh",
  },
  {
    label: "English",
    value: "en",
  },
];
const items = [i18n("language")];

const activeIndex = ref(0);

const modelvalue = ref(currentLang);

const selectItem = (index: number) => {
  activeIndex.value = index;
};

async function submit() {
  setLang(modelvalue.value);
  Dialog.showMessageBox({
    message: i18n("save.success"),
    title: i18n("language"),
    type: "info",
  }).then(() => {
    location.reload();
  });
}
</script>
<style scoped>
@import "./setStyle.css";
</style>
