import { GetClientId } from "@/util/clientid";
import { getSystemConfig, setSystemConfig } from "./config";
const config = getSystemConfig();
export async function authWithWechat(): Promise<boolean> {
  const queryParams = new URLSearchParams(window.location.search);
  const code = queryParams.get("code");
  if (!code) {
    alert("登录失败，无法继续操作");
    return false;
  }
  return await toLogin(code);
}


// 登录接口
async function toLogin(code: string): Promise<boolean> {
  const data = {
    login_type: "qyweixin",
    client_id: GetClientId(),
    param: {
      code,
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
      window.location.href = "/";
      return true;
    }
  } catch (error) {
    alert(error);
  }
  return false;
}