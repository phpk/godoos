import axios, { AxiosHeaders } from 'axios';
import { errMsg } from './msg';
import { getClientId } from './uuid'
import { API_URL } from '../stores/config';
const baseURL = '/api/'
/** 创建axios实例 */
const isDev = typeof process !== 'undefined' && process?.env?.NODE_ENV === 'development';
const instanceBaseURL = isDev ? baseURL : `${API_URL}${baseURL}`;

const request = axios.create({
  baseURL: instanceBaseURL,
  timeout: 60000
})
export const getToken = () => {
  // 如果 URL 中没有 token，则从 localStorage 中获取
  return localStorage.getItem("Authorization");
}
export const setToken = (token: string) => {
  localStorage.setItem("Authorization", token);
}
export const getUsername = () => {
  return localStorage.getItem("GodoOS-username") || "admin";
}
export const setUsername = (username: string) => {
  localStorage.setItem("GodoOS-username", username);
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
    config.headers = config.headers || new AxiosHeaders()
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
      //console.log(err)
      if (err.response) {
        const status = err.response.status;
        const message = err.response.data;
        
        if (status === 404) {
          errMsg('404错误: 请求的资源未找到');
        } else if (status === 500) {
          errMsg('500错误: 服务器内部错误');
        } else if (err.code === 'ERR_NETWORK') {
          errMsg('网络连接错误,请稍后重试');
        } else if (err.code === 'ECONNABORTED') {
          errMsg('服务异常,请稍后重试');
        } else {
          errMsg('未知错误: ' + message);
        }
        //return;
        return Promise.reject(false);
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