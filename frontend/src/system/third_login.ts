import { GetClientId } from "@/util/clientid";

import { useLoginStore } from "@/stores/login";
import { notifyError } from "@/util/msg";
import { getSystemConfig, setSystemConfig } from "./config";
const config = getSystemConfig();
export async function authWithThirdParty(unionid: string): Promise<boolean> {
  const data = {
    login_type: "third_api",
    client_id: GetClientId(),
    param: {
      unionid,
    },
  };
  try {
    const res = await fetch("http://server001.godoos.com/user/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
    const jsondata = await res.json();
    if (jsondata.success) {
      jsondata.data.url = "http://server001.godoos.com";
      config.userInfo = jsondata.data;
      config.userType = "member";
      setSystemConfig(config);

      if (jsondata.code == 10000) {
        const loginStore = useLoginStore();
        loginStore.ThirdPartyLoginMethod = "register";

        // 注册信息
        loginStore.registerInfo = {
          username: jsondata.data.username,
          nickname: jsondata.data.nickname,
          password: jsondata.data.password,
          email: jsondata.data.email,
          phone: jsondata.data.phone,
          third_user_id: jsondata.data.third_user_id,
          union_id: jsondata.data.union_id,
          patform: jsondata.data.patform,
          confirmPassword: "",
        }
        return true;
      }

      // config.userInfo.username = "";
      // config.userInfo.password = "";
      // setSystemConfig(config);

      window.location.href = "/";
      return true;
    }
    notifyError("登录失败");
    setTimeout(() => {
      window.location.href = "/";
    }, 5000);
  } catch (error) {
    notifyError("登录失败");
    console.error(error);
  }
  return false;
}