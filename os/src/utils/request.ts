import axios from 'axios';
import { errMsg } from './msg';
import { getClientId } from './uuid'
const baseURL = '/'
/** 创建axios实例 */
const request = axios.create({
  baseURL,
  timeout: 60000
})
export const getToken = () => {
  // 如果 URL 中没有 token，则从 localStorage 中获取
  return localStorage.getItem("Authorization");
}
export const setToken = (token: string) => {
  localStorage.setItem("Authorization", token);
}
export function getHeader() {
  return {
    'Content-Type': 'application/json',
    'Authorization': getToken(),
    'ClientID': getClientId(),
  }
}
export const clearToken = () => {
  localStorage.removeItem("Authorization");
  localStorage.removeItem("ClientID");
}
export function getUrl(url: string, islast = true) {
  if (islast) {
    return url + '&uuid=' + getClientId() + '&token=' + getToken()
  } else {
    return url + '?uuid=' + getClientId() + '&token=' + getToken()
  }
}
// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 发请求带上token
    const token = getToken()
    config.headers = config.headers || {}
    const clientId = getClientId()
    if (clientId) {
      config.headers['ClientId'] = clientId
    }
    if (token) {
      config.headers['Authorization'] = token
    } else {
      // 如果没有 token，则跳转到登录页面
      // router.push('/login')
      return Promise.reject(new Error('No token found, redirecting to login'))
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (resp: any) => {
    // 响应成功的拦截
    if (resp) {
      //console.log(resp)
      if (resp.data && resp.data.code === -2) {
        errMsg('登录已过期,请重新登录')
      }
      return resp.data
    }
  },
  (err) => {
    // 处理错误信息
    if (err) {
      console.log(err)
      if (err.code === 'ERR_NETWORK') {
        errMsg('网络连接错误,请稍后重试')
        return Promise.reject(err.response)
      } else if (err.code === 'ECONNABORTED') {
        errMsg('服务异常,请稍后重试')
        return Promise.reject(err.response)
      }
      if (err.response) {
        const message = err.response.data
        console.log('message:', message)
        errMsg(message)
        return Promise.reject(err.response.data)
      }

    }
  }
)

interface RequestParams {
  [key: string]: any
}

// 定义基础响应格式
interface BaseResponse {
  code: number
  message: string
  success: boolean
  time: number
}

// 将基础响应和具体数据结合的泛型类型
export type ApiResponse<T> = BaseResponse & {
  data: T
}

/**
 * 通用GET请求方法
 * @param url 请求地址
 * @param params 请求参数
 * @returns Promise<ApiResponse<T>>
 */
export const get = <T = any>(
  url: string,
  params?: RequestParams
): Promise<ApiResponse<T>> => {
  return request({
    method: 'GET',
    url,
    params,
  })
}

/**
 * 通用POST请求方法
 * @param url 请求地址
 * @param data 请求数据
 * @param params 查询参数（可选）
 * @returns Promise<ApiResponse<T>>
 */
export const post = <T = any>(
  url: string,
  data?: RequestParams,
  params?: RequestParams
): Promise<ApiResponse<T>> => {
  return request({
    method: 'POST',
    url,
    data,
    params,
  })
}