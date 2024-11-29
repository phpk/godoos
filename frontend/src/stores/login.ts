import { defineStore } from 'pinia';

export const useLoginStore = defineStore("login", () => {
  // 第三方登录方式
  const ThirdPartyPlatform = ref<string | null>(null);

  const State = ref<string | null>(null);
  return {
    ThirdPartyPlatform,
    State
  };
});