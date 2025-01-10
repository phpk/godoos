import { defineStore } from 'pinia';

export const useLoginStore = defineStore("login", () => {
  // 第三方登录方式
  const ThirdPartyPlatform = ref<string | null>(null);

  const State = ref<string>("");
  // 临时token
  const tempToken = ref<string>("");
  // 临时clientId
  const tempClientId = ref<string>("");
  const ThirdPartyLoginMethod = ref("login");

  const thirdpartyList = ref([]);

  const registerInfo = ref({
    "username": "",
    "nickname": "",
    "password": "",
    "email": "",
    "phone": "",
    "third_user_id": "",
    "union_id": "",
    "patform": "",
    "confirmPassword": "",
  });
  return {
    ThirdPartyPlatform,
    State,
    ThirdPartyLoginMethod,
    tempToken,
    tempClientId,
    registerInfo,
    thirdpartyList,
  };
},);