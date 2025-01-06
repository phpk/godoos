import { defineStore } from 'pinia';

export const useLoginStore = defineStore("login", () => {
  // 第三方登录方式
  const ThirdPartyPlatform = ref<string | null>(null);

  const State = ref<string>("");
  // 临时token
  const tempToken = ref<string>("");
  // 临时clientId
  const tempClientId = ref<string>("");
  const page = ref("login");

  return {
    ThirdPartyPlatform,
    State,
    page,
    tempToken,
    tempClientId,
  };
},);