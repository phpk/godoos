import { defineStore } from "pinia";
import { ref } from "vue";
import { systemSettingList, settingsConfig } from "@/stores/settingsConfig";
import { md5 } from "js-md5";
import { type as osType, platform } from "@tauri-apps/plugin-os";
export const useSettingsStore = defineStore('settings', () => {
  const settingList = ref(systemSettingList);

  const config: any = ref(settingsConfig);
  const isLockScreen = ref(false);
  const systemInfo: any = ref({});
  const updateConfig = (newConfig: any) => {
    config.value = newConfig;
  }
  const setConfig = (key: string, value: any) => {
    config.value[key] = value;
  }
  const getConfig = (key: string) => {
    return config.value[key]
  }
  const initConfig = () => {
    config.value = settingsConfig;
  }

  const setLockTime = () => {
    if (config.value.lock.timeout > 0 && !isLockScreen.value) {
      config.value.lock.activeTime = new Date().getTime() + config.value.lock.timeout * 1000;
    }
  }
  const checkLockPassword = (password: string) => {
    return config.value.lock.password === md5(password);
  }
  const unLockScreen = () => {
    config.value.lock.activeTime = new Date().getTime() + config.value.lock.timeout * 1000;
    isLockScreen.value = false;
  }
  const checkIsLock = () => {
    if (config.value.lock.timeout > 0) {
      const now = new Date().getTime();
      isLockScreen.value = now > config.value.lock.activeTime;
      //console.log(isLockScreen.value)
      return isLockScreen.value;
    }
    isLockScreen.value = false;
    return false;
  }
  const osInit = () => {
    try {
      systemInfo.value.platform = platform();
    }
    catch (error) {
      systemInfo.value.platform = "web";
    }
    try {
      systemInfo.value.ostype = osType();
    }
    catch (error) {
      // console.warn(error);
      systemInfo.value.ostype = "web";
    }
    systemInfo.value.isMobile = ["android", "ios"].includes(systemInfo.value.ostype);
    systemInfo.value.isWeb = systemInfo.value.platform === "web";
    systemInfo.value.isDesktop = !systemInfo.value.isMobile && !systemInfo.value.isWeb;
  }
  return {
    settingList,
    config,
    isLockScreen,
    updateConfig,
    setConfig,
    getConfig,
    initConfig,
    setLockTime,
    checkIsLock,
    unLockScreen,
    checkLockPassword,
    osInit,
    systemInfo,
  }
}, {
  persist: {
    key: 'settingsStore',
    pick: ['config', 'systemInfo'],
  },
})
