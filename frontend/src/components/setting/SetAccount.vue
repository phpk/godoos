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
          <h1 class="setting-title">{{ i18n('account.info') }}</h1>
        </div>
        <div class="setting-item">
          <label> {{ i18n('account') }} </label>
          <el-input v-model="account.username" :placeholder="i18n('account')" clearable />
        </div>
        <div class="setting-item">
          <label> {{ i18n('password') }} </label>
          <el-input v-model="account.password" type="password" :placeholder="i18n('password')" clearable/>
        </div>

        <div class="setting-item">
          <label></label>
          <el-button @click="submit" type="primary">
            {{ i18n('confirm') }} 
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
import { Dialog,i18n } from '@/system/index.ts';
import { getSystemKey, setSystemKey } from '@/system/config'

const items = [i18n('account.info')];

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
    message: i18n('save.success'),
    title: i18n('account'),
    type: 'info',
  }).then(() => {
    location.reload();
  });
}
</script>
<style scoped>
@import './setStyle.css';
</style>
