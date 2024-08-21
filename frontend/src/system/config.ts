import { generateRandomString } from "../util/common.ts"
export const configStoreType = localStorage.getItem('GodoOS-storeType') || 'browser';
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
        config.version = '1.0.0';
        //config.version = '0.0.9';
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

    // 初始化用户信息，若本地存储中已存在则不进行覆盖
    if (!config.userInfo) {
        config.userInfo = {
            serverUrl:'',
            username: '',
            password: '',
            memberId: 0,
            nickname: '',
            avatar: '',
            email: '',
            mobile: '',
            role: '',
            department: ''
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
            memberId: '',
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
    if(config.storeType == 'net'){
        return config.storenet.url + '/file'
    }
    else if (config.storeType == 'webdav') {
        return config.apiUrl + '/webdav'
    }else{
        return config.apiUrl + '/file'
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