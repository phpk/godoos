import { GetClientId } from "../util/clientid.ts";
import { generateRandomString } from "../util/common.ts";
import { parseAiConfig } from "./aiconfig.ts";
export const configStoreType = localStorage.getItem('GodoOS-storeType') || 'local';
/**
 * 获取系统配置信息。
 * 从本地存储中获取或初始化系统配置对象，并根据条件决定是否更新本地存储中的配置。
 * @param ifset 是否将配置信息更新回本地存储
 * @returns 当前系统配置对象
 */
export const getSystemConfig = (ifset = false) => {
  // 从本地存储中尝试获取配置信息，若不存在则使用默认空对象
  const configSetting = localStorage.getItem('GodoOS-config') || '{}';
  // 解析配置信息为JSON对象
  let config: any = JSON.parse(configSetting);

  // 初始化配置对象的各项属性，若本地存储中已存在则不进行覆盖
  if (!config.version) {
    config.version = '1.0.6';
  }
  if (!config.isFirstRun) {
    config.isFirstRun = false;
  }
  if (!config.lang) {
    config.lang = '';
  }
  // 初始化API相关URL，若本地存储中已存在则不进行覆盖
  if (!config.apiUrl) {
    config.apiUrl = `${window.location.protocol}//${window.location.hostname}:56780`;
  }
  if (!config.userType) {
    config.userType = 'person'//如果是企业版请改为member
  }

  if (!config.editorType) {
    config.editorType = 'local'
  }
  if (!config.onlyoffice) {
    config.onlyoffice = {
      url: '',
      sceret: '',
    }
  }
  if (!config.file) {
    config.file = {
      isPwd: 0,
      pwd: '',
      // salt: 'vIf_wIUedciAd0nTm6qjJA=='
    }
  }
  if (!config.fileInputPwd) {
    config.fileInputPwd = []
  }
  // 初始化用户信息，若本地存储中已存在则不进行覆盖
  if (!config.userInfo) {
    config.userInfo = {
      url: '',//如果是企业版请改为服务器地址，不要加斜线
      username: '',
      password: '',
      id: 0,
      nickname: '',
      avatar: '',
      email: '',
      phone: '',
      desc: '',
      job_number: '',
      work_place: '',
      hired_date: '',
      ding_id: '',
      role_id: 0,
      roleName: '',
      dept_id: 0,
      deptName: '',
      token: '',
      user_auths: '',
      user_shares: '',
      isPwd: false
    };
  }

  config.isApp = (window as any).go ? true : false;
  // 初始化系统相关信息，若本地存储中已存在则不进行覆盖
  //system
  if (!config.systemInfo) {
    config.systemInfo = {};
  }
  if (!config.theme) {
    config.theme = 'light';
  }
  if (!config.storeType) {
    config.storeType = configStoreType;
  }
  if (!config.storePath) {
    config.storePath = "";
  }
  if (!config.netPath) {
    config.netPath = "";
  }
  if (!config.netPort) {
    config.netPort = "56780";
  }
  // 初始化背景设置，若本地存储中已存在则不进行覆盖
  if (!config.background) {
    config.background = {
      url: '/image/bg/bg6.jpg',
      type: 'image',
      color: 'rgba(30, 144, 255, 1)',
      imageList: [
        '/image/bg/bg1.jpg',
        '/image/bg/bg2.jpg',
        '/image/bg/bg3.jpg',
        '/image/bg/bg4.jpg',
        '/image/bg/bg5.jpg',
        '/image/bg/bg6.jpg',
        '/image/bg/bg7.jpg',
        '/image/bg/bg8.jpg',
        '/image/bg/bg9.jpg',
      ]
    }
  }

  // 初始化账户信息，若本地存储中已存在则不进行覆盖
  if (!config.account) {
    config.account = {
      ad: false,
      username: '',
      password: '',
    };
  }
  if (config.userType == 'member') {
    config.account.ad = false
  }
  if (!config.storenet) {
    config.storenet = {
      url: '',
      username: '',
      password: '',
      isCors: ''
    };
  }
  if (!config.webdavClient) {
    config.webdavClient = {
      url: '',
      username: '',
      password: '',
    };
  }
  if (!config.dbInfo) {
    config.dbInfo = {
      url: '',
      username: '',
      password: '',
      dbname: ''
    };
  }
  if (!config.chatConf) {
    config.chatConf = {
      'checkTime': '15',
      'first': '192',
      'second': '168',
      'thirdStart': '1',
      'thirdEnd': '1',
      'fourthStart': '2',
      'fourthEnd': '254'
    }
  }

  // 初始化桌面快捷方式列表，若本地存储中已存在则不进行覆盖
  if (!config.desktopList) {
    config.desktopList = [];
  }
  // 初始化菜单列表，若本地存储中已存在则不进行覆盖
  if (!config.menuList) {
    config.menuList = [];
  }
  // 生成新的会话ID，若本地存储中不存在
  if (!config.token) {
    config.token = generateRandomString(16);
  }
  config = parseAiConfig(config);
  // 根据参数决定是否更新本地存储中的配置信息
  if (ifset) {
    setSystemConfig(config)
  }
  // 返回配置对象
  return config;
};

export function getApiUrl() {
  const config = getSystemConfig();
  if (config.userType == 'person') {
    return config.apiUrl
  } else {
    return config.userInfo.url
  }
}
export function getFileUrl() {
  const config = getSystemConfig();
  if (config.userType == 'person') {
    if (config.storeType == 'net') {
      return config.storenet.url + '/file'
    }
    else if (config.storeType == 'webdav') {
      return config.apiUrl + '/webdav'
    } else {
      return config.apiUrl + '/file'
    }
  } else {
    return config.userInfo.url + '/files'
  }

}
export function getChatUrl() {
  const config = getSystemConfig();
  if (config.userType == 'person') {
    return config.apiUrl + '/localchat'
  } else {
    return config.userInfo.url + '/chat'
  }
}
export function getUrl(url: string, islast = true) {
  const config = getSystemConfig();
  if (config.userType == 'person') {
    return config.apiUrl + url
  } else {
    if (islast) {
      return config.userInfo.url + url + '&uuid=' + GetClientId() + '&token=' + config.userInfo.token
    } else {
      return config.userInfo.url + url + '?uuid=' + GetClientId() + '&token=' + config.userInfo.token
    }

  }
}
export function getWorkflowUrl() {
  const config = getSystemConfig();
  if (config.userType == 'member') {
    console.log(config.userInfo)
    return config.userInfo.url + '/views/desktop/index.html' + '?uuid=' + GetClientId() + '&token=' + config.userInfo.token
  }
}
export function fetchGet(url: string, headerConfig?: { [key: string]: string }) {
  // console.log('请求头部；', headerConfig);

  const config = getSystemConfig();
  if (config.userType == 'person') {
    return fetch(url, {
      method: 'GET',
      headers: headerConfig
    })
  } else {
    return fetch(url, {
      method: 'GET',
      credentials: 'include',
      headers: {
        //'Content-Type': 'application/json',
        'ClientID': GetClientId(),
        'Authorization': config.userInfo.token,
        ...headerConfig
      }
    })
  }
}
export function fetchPost(url: string, data: any, headerConfig?: { [key: string]: string }) {
  const config = getSystemConfig();
  if (config.userType == 'person') {
    return fetch(url, {
      method: 'POST',
      body: data,
      headers: headerConfig
    })
  } else {
    return fetch(url, {
      method: 'POST',
      credentials: 'include',
      body: data,
      headers: {
        'ClientID': GetClientId(),
        'Authorization': config.userInfo.token,
        ...headerConfig
      }
    })
  }
}
export function isWindowsOS() {
  return /win64|wow64|win32|win16|wow32/i.test(navigator.userAgent);
}
export function parseJson(str: string) {
  try {
    return JSON.parse(str);
  } catch (e) {
    return undefined;
  }
};
export function getSplit() {
  if (isWindowsOS()) {
    return "\\"
  } else {
    return "/"
  }
}

export const getSystemKey = (key: string, ifset = false) => {
  const config = getSystemConfig(ifset);
  if (key.indexOf('.') > -1) {
    const keys = key.split('.');
    return config[keys[0]][keys[1]];
  } else {
    return config[key];
  }
}

export const setSystemKey = (key: string, val: any) => {
  const config = getSystemConfig();
  config[key] = val;
  localStorage.setItem('GodoOS-config', JSON.stringify(config));
  localStorage.setItem('GodoOS-storeType', config.storeType);
};

export const setSystemConfig = (config: any) => {
  localStorage.setItem('GodoOS-config', JSON.stringify(config));
  localStorage.setItem('GodoOS-storeType', config.storeType);
};

export const clearSystemConfig = () => {
  const config = getSystemConfig();

  if (config.userType === 'person') {
    const storetype = localStorage.getItem('GodoOS-storeType') || 'local';
    localStorage.clear()
    localStorage.setItem('GodoOS-storeType', storetype)
  }
  //localStorage.removeItem('GodoOS-config');
};
