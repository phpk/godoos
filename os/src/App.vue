<!-- src/App.vue -->
<template>
  <el-config-provider :locale="currentLocale">
    <WinScreen />
  </el-config-provider>
</template>

<script setup lang="ts">
import { ElConfigProvider } from 'element-plus';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import en from 'element-plus/es/locale/lang/en';
import { getLang, setLang } from '@/i18n';
import { computed, ref, onMounted, onUnmounted } from 'vue';
import { eventBus } from '@/interfaces/event';

// 获取当前语言
const currentLang = ref(getLang());

// 动态更新 ElConfigProvider 的 locale 属性
const currentLocale = computed(() => {
  return currentLang.value === 'en' ? en : zhCn;
});

const setLanguage = (lang: string) => {
  currentLang.value = lang;
  setLang(lang);
};

// 监听 setLanguages 事件
const onSetLanguages = (lang: string) => {
  setLanguage(lang);
};

onMounted(() => {
  eventBus.on('setLanguages', onSetLanguages);
});

onUnmounted(() => {
  eventBus.off('setLanguages', onSetLanguages);
});
</script>