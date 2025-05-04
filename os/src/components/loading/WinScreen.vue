<template>
  <div class="loading" v-if="!isLoading">
    <WinLogo />
    <div id="text"></div>
    <WinLoading />
  </div>
  <template v-else>
    <LockScreen v-if="loginStore.isLoginState && settingsStore.isLockScreen" />
    <template v-else>
      <div class="desktop" v-if="loginStore.isLoginState" ref="desktopRef" :style="desktopBackgroundStyle">
        <template v-if="!isMobileDevice()">
          <desktop-icons :isDesktop="true" :fileList="desktopStore.icons" showClass="desktop-icons" />
          <FileContextMenu />
          <ComputerContextMenu />
          <window-manager />
          <taskbar />

        </template>
        <template v-else>
          <desktop-mobile :isDesktop="true" :fileList="desktopStore.mobileicons" />
          <FileContextMenu />
          <window-manager />
        </template>
      </div>
      <AuthLogin v-else />
    </template>
  </template>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useLoginStore } from "@/stores/login";
import { useDesktopStore } from '@/stores/desktop';
import { useSettingsStore } from "@/stores/settings";
import { useMessageStore } from '@/stores/message';
import { useAssistantStore } from '@/stores/assistant';
import { useModelStore } from '@/stores/model';
import { isMobileDevice } from '@/utils/device';
import LockScreen from './LockScreen.vue';

const loginStore = useLoginStore();
const desktopStore = useDesktopStore();
const settingsStore = useSettingsStore();
const messageStore = useMessageStore();
const assistantStore = useAssistantStore();
const modelStore = useModelStore();
const isLoading = ref(false);
const desktopRef = ref<HTMLDivElement | null>(null);

onMounted(async () => {
  await desktopStore.initDesktop();
  isLoading.value = true;
  setTimeout(() => {
    initSystem();
  }, 1000);
});
const initSystem = async () => {
  await settingsStore.osInit();
  messageStore.initMessage();
  await assistantStore.initPrompt();
  await modelStore.initModel();
};

watch(
  () => loginStore.isLoginState,
  () => {
    if (loginStore.isLoginState) {
      //setDesktopBackground();// 登录成功后，将焦点设置到桌面元素上
      messageStore.initMessage();
    } else {
      messageStore.closeMessage();
    }
  }
);
// 计算属性来动态设置背景样式
const desktopBackgroundStyle = computed(() => {
  const background = settingsStore.config.background;
  if (background.type === "image") {
    return {
      backgroundImage: `url(${background.url})`,
      backgroundColor: ""
    };
  } else if (background.type === "color") {
    return {
      backgroundColor: background.color,
      backgroundImage: ""
    };
  }
  return {};
});
</script>

<style lang="scss" scoped>
@use "@/styles/screen.scss";
</style>