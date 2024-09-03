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
      <div v-if="1 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">广告与更新提示</h1>

        </div>
        <div class="setting-item">
          <label></label>
          <el-switch 
          v-model="ad" 
          active-text="开启" 
          inactive-text="关闭" 
          size="large"
            :before-change="setAd"></el-switch>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>

import { ref } from 'vue';
import { Dialog, t } from '@/system/index.ts';
import { getSystemKey, setSystemKey } from '@/system/config'
import { ElMessageBox } from 'element-plus'
const items = ['锁屏设置', '广告设置'];
const activeIndex = ref(0);
const account = ref(getSystemKey('account'));
const ad = ref(account.value.ad)

const selectItem = (index: number) => {
  activeIndex.value = index;
};

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
  }else{
    data.ad = true
    setSystemKey('account', data);
    return Promise.resolve(true)
  }

}
</script>
<style scoped>
@import './setStyle.css';
</style>
