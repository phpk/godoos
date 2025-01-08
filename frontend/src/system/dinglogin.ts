import { GetClientId } from "@/util/clientid";
import dd from "dingtalk-jsapi";
import { getSystemConfig, setSystemConfig } from "./config";

const config = getSystemConfig();
// 函数用于动态加载外部JS文件
export function loadScript(url: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = url
    script.onload = () => resolve()
    script.onerror = () => reject(new Error(`Failed to load script ${url}`))
    document.head.appendChild(script)
  })
}

const currentUrl = window.location.origin

export async function authWithDing(): Promise<boolean> {
  try {
    const res = await fetch(currentUrl + "/user/ding/conf");
    const data = await res.json();
    console.log(data)
    if (data.success) {
      console.log(data.data.id)
      return await getCode(data.data.id);
    }
    return false;
  } catch (e) {
    console.error(e);
    return false;
  }
}

// 获取code登录
async function getCode(corpId: string): Promise<boolean> {
  // 加载钉钉登录脚本
  await loadScript("https://g.alicdn.com/dingding/dingtalk-jsapi/3.0.12/dingtalk.open.js")

  return new Promise((resolve) => {
    dd.runtime.permission.requestAuthCode({
      corpId: corpId,
      onSuccess: async function (result: { code: string }) {
        resolve(await toLogin(result.code));
      },
      onFail: function (err: any) {
        console.log(err)
        resolve(false);
      },
    } as any);
  });
}


// 登录接口
async function toLogin(code: string): Promise<boolean> {
  const data = {
    login_type: "dingtalk_workbench",
    client_id: GetClientId(),
    param: {
      code,
    },
  };
  const res = await fetch(currentUrl + "/user/login", {
    method: "POST",
    body: JSON.stringify(data),
  });
  const jsondata = await res.json();
  if (jsondata.success) {
    jsondata.data.url = currentUrl;
    config.userInfo = jsondata.data;
    config.userType = "member";
    setSystemConfig(config);
    return true;
  }
  return false;
}
