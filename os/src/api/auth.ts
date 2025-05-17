import { loadScript } from '@/utils/load'
import { get, getToken, setToken,post } from '@/utils/request'
import { getClientId } from '@/utils/uuid'
import { errMsg } from '@/utils/msg';
export function loginIn(params: any) {
  return post('user/login', params).then(res => {
    if (res.success) {
      setToken(res.data.token)
    }
    return res
  }).catch(err => { 
    errMsg("请求失败")
    throw new Error(err) 
  })
}
export async function logout() {
  return get('user/logout')
}
export async function isLogin() {
  const res = await get('user/islogin')
  return res.success
}
export async function getDingConf() {
  await loadScript(
    "https://g.alicdn.com/dingding/h5-dingtalk-login/0.21.0/ddlogin.js"
  );

  return await get("user/ding/conf");
}
export async function getThirdpartyList() {
  const result = await get("/user/thirdparty/list"
  );
  if (result.success) {
    return result.data.list;
  }
  return [];
};
export async function getEmailCode(email: string) {
  const data = {
    email: email,
    client_id: getClientId(),
  }
  return await post('/user/emailcode', data)
}
export async function getSmsCode(phone: string) {
  const data = {
    phone: phone,
    client_id: getClientId(),
  }
  return await post('/user/smscode', data)
}
export async function register(params: any) {
  console.log(params)
  return post('/user/register', params)
}