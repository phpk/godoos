import { generateRandomString } from "../util/common.ts"
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
    const config = JSON.parse(configSetting);

    // 初始化配置对象的各项属性，若本地存储中已存在则不进行覆盖
    if (!config.version) {
        config.version = '1.0.1';
    }
    if (!config.isFirstRun) {
        config.isFirstRun = false;
    }
    if (!config.lang) {
        config.lang = '';
    }
    // 初始化API相关URL，若本地存储中已存在则不进行覆盖
    if (!config.apiUrl) {
        config.apiUrl = 'http://localhost:56780';
    }
    if (!config.userType) {
        config.userType = 'person'
    }
    // 初始化用户信息，若本地存储中已存在则不进行覆盖
    if (!config.userInfo) {
        config.userInfo = {
            url: '',
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
            token: ''
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
    if (!config.userType) {
        config.userType = 'person';
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
            username: '',
            password: '',
        };
    }
    if (!config.storenet) {
        config.storenet = {
            url: '',
            username: '',
            password: '',
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
    // 根据参数决定是否更新本地存储中的配置信息
    if (ifset) {
        setSystemConfig(config)
    }
    // 返回配置对象
    return config;
};

export function getApiUrl() {
    return getSystemKey('apiUrl')
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
    }else{
        return config.userInfo.url + '/files'
    }

}
export function fetchGet(url: string) {
    const config = getSystemConfig();
    if (config.userType == 'person') {
        return fetch(url)
    } else {
        return fetch(url, {
            method: 'GET',
            credentials: 'include',
            headers: {
                //'Content-Type': 'application/json',
                'ClientID': getClientId(),
                'Authorization': config.userInfo.token
            }
        })
    }
}
export function fetchPost(url: string, data: any) {
    const config = getSystemConfig();
    if (config.userType == 'person') {
        return fetch(url, {
            method: 'POST',
            body: data,
        })
    } else {
        return fetch(url, {
            method: 'POST',
            credentials: 'include',
            body: data,
            headers: {
                'ClientID': getClientId(),
                'Authorization': config.userInfo.token
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
        const storeType = getSystemKey('storeType')
        //console.log(storeType)
        if (storeType === 'browser') {
            return "/"
        } else {
            return "\\"
        }
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
    const storetype = localStorage.getItem('GodoOS-storeType') || 'browser';
    localStorage.clear()
    localStorage.setItem('GodoOS-storeType', storetype)
    //localStorage.removeItem('GodoOS-config');
};
function bin2hex(s:string) {
    s = encodeURI(s);//只会有0-127的ascii不转化
    let m:any = s.match(/%[\dA-F]{2}/g), a:any = s.split(/%[\dA-F]{2}/), i, j, n, t;
    m.push("")
    for (i in a) {
      if (a[i] === "") { a[i] = m[i]; continue }
      n = ""
      for (j in a[i]) {
        t = a[i][j].charCodeAt().toString(16).toUpperCase()
        if (t.length === 1) t = "0" + t
        n += "%" + t
      }
      a[i] = n + m[i]
    }
    return a.join("").split("%").join("")
  }
  export const getClientId = () => {
    let uuid:any = localStorage.getItem("godoosClientId");
    if (!uuid) {
      let canvas = document.createElement('canvas');
      let ctx:any = canvas.getContext('2d');
      ctx.fillStyle = '#FF0000';
      ctx.fillRect(0, 0, 8, 10);
      let b64 = canvas.toDataURL().replace("data:image/png;base64,", "");
      let bin = window.atob(b64);
      uuid = bin2hex(bin.slice(-16, -12));
      localStorage.setItem("godoosClientId", uuid);
    }
    return uuid;
  }