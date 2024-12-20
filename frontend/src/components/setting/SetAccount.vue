<template>
  <div class="container">
    <div class="nav">
      <ul>
        <li v-for="(item, index) in items" :key="index" @click="selectItem(index)"
          :class="{ active: index === activeIndex }">
          {{ item }}
        </li>
      </ul>
    </div>
    <div class="setting">
      <div v-if="0 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">{{ t("background") }}</h1>
        </div>
        <div class="setting-item">
          <el-select v-model="config.background.type">
            <el-option v-for="(item, key) in desktopOptions" :key="key" :label="item.label" :value="item.value" />
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
              <li v-for="(item, index) in config.background.imageList" :key="index"
                :class="config.background.url === item ? 'selected' : ''" @click="setBg(item)">
                <img :src="item" />
              </li>
            </ul>
          </div>
          <div class="setting-item">
            <label> </label>
          </div>
        </template>
      </div>
      <div v-if="1 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">锁屏</h1>
        </div>
        <div class="setting-item">
          <label> {{ t('account') }} </label>
          <el-input v-model="account.username" :placeholder="t('account')" clearable />
        </div>
        <div class="setting-item">
          <label> {{ t('password') }} </label>
          <el-input v-model="account.password" type="password" :placeholder="t('password')" clearable />
        </div>

        <div class="setting-item">
          <label></label>
          <el-button @click="submit" type="primary">
            {{ t('confirm') }}
          </el-button>
        </div>
      </div>
      <div v-if="2 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">广告与更新提示</h1>

        </div>
        <div class="setting-item">
          <label></label>
          <el-switch v-model="ad" active-text="开启" inactive-text="关闭" size="large" :before-change="setAd"></el-switch>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>

import { ref } from 'vue';
import { Dialog, t, useSystem } from '@/system/index.ts';
import { getSystemKey, setSystemKey, getSystemConfig, setSystemConfig } from '@/system/config'
import { ElMessageBox } from 'element-plus'
const sys = useSystem();
const items = [t("background"), '锁屏设置', '广告设置'];
const activeIndex = ref(0);
const account = ref(getSystemKey('account'));
const ad = ref(account.value.ad)
const config: any = ref(getSystemConfig());
const selectItem = (index: number) => {
  activeIndex.value = index;
};
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
async function submit() {
  setSystemKey('account', account.value);
  Dialog.showMessageBox({
    message: t('save.success'),
    title: t('account'),
    type: 'info',
  }).then(() => {
    location.reload();
  });
}
function setAd() {
  const data = toRaw(account.value)
  if (ad.value) {
    return new Promise((resolve) => {
      setTimeout(() => {
        ElMessageBox.confirm(
          '广告关闭后您将收不到任何系统通知和更新提示！',
          'Warning',
          {
            confirmButtonText: '确定关闭',
            cancelButtonText: '取消',
            type: 'warning',
          }
        )
          .then(() => {
            data.ad = false
            setSystemKey('account', data);
            return resolve(true)
          })
          .catch(() => {
            return resolve(false)
          })

      }, 1000)
    })
  } else {
    data.ad = true
    setSystemKey('account', data);
    return Promise.resolve(true)
  }

}
function setBg(item: any) {
  config.value.background.url = item
  config.value.background.type = "image";
  setSystemConfig(config.value);
  sys.initBackground();

}
function onColorChange(color: string) {
  config.value.background.color = color;
  config.value.background.type = "color";
  setSystemConfig(config.value);
  sys.initBackground();
}

</script>
<style scoped>
@import "./setStyle.css";
@import "@/assets/imglist.scss";
</style>
