import { loadScript } from '@/utils/load'
import { get, getToken, setToken,post } from '@/utils/request'
import { getClientId } from '@/utils/uuid'

export function loginIn(params: any) {
  return fetch('/user/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(params)
  }).then(res => res.json()).then(res => {
    if (res.success) {
      setToken(res.data.token)
    }
    return res
  }).catch(err => { throw new Error(err) })
}
export async function logout() {
  return post('user/logout', {
    method: 'POST',
  })
}
export async function isLogin() {
  const token = getToken()
  if (!token) {
    return false
  }
  const res = await get('user/islogin')
  return res.success
}
export async function getDingConf() {
  await loadScript(
    "https://g.alicdn.com/dingding/h5-dingtalk-login/0.21.0/ddlogin.js"
  );

  const res = await fetch("user/ding/conf");
  return await res.json();
}
export async function getThirdpartyList() {
  const result = await fetch("/user/thirdparty/list"
  );
  if (result.ok) {
    const data = await result.json();
    if (data.success) return data.data.list;
  }
  return [];
};
export async function getEmailCode(email: string) {
  const data = {
    email: email,
    client_id: getClientId(),
  }
  const res = await fetch('/user/emailcode', {
    method: 'POST',
    body: JSON.stringify(data),
  })
  return await res.json()
}
export async function getSmsCode(phone: string) {
  const data = {
    phone: phone,
    client_id: getClientId(),
  }
  const res = await fetch('/user/smscode', {
    method: 'POST',
    body: JSON.stringify(data),
  })
  return await res.json()
}
export async function register(params: any) {
  return fetch('/user/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(params)
  }).then(res => res.json())
}