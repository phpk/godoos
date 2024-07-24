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
          <h1 class="setting-title">{{ t('account.info') }}</h1>
        </div>
        <div class="setting-item">
          <label> {{ t('account') }} </label>
          <el-input v-model="account.username" :placeholder="t('account')" clearable />
        </div>
        <div class="setting-item">
          <label> {{ t('password') }} </label>
          <el-input v-model="account.password" type="password" :placeholder="t('password')" clearable/>
        </div>

        <div class="setting-item">
          <label></label>
          <el-button @click="submit" type="primary">
            {{ t('confirm') }} 
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
// import WinButton from '@/components/win/WinButton.vue';
// import WinSelect from '@/components/win/WinSelect.vue';

import { ref } from 'vue';
import { Dialog,t } from '@/system/index.ts';
import { getSystemKey, setSystemKey } from '@/system/config'

const items = [t('account.info')];

const activeIndex = ref(0);

// const modelvalue = ref(system.getConfig('lang'));
// const password = ref('');
const account = ref(getSystemKey('account'));

const selectItem = (index: number) => {
  activeIndex.value = index;
};

async function submit() {
  //   await system.setConfig('lang', modelvalue.value);
  // localStorage.setItem('godoOS_username', account.value);
  // localStorage.setItem('godoOS_password', password.value);
  setSystemKey('account', account.value);
  Dialog.showMessageBox({
    message: t('save.success'),
    title: t('account'),
    type: 'info',
  }).then(() => {
    location.reload();
  });
}
</script>
<style scoped>
@import './setStyle.css';
</style>
