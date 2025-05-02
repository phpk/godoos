import bg1 from "/public/images/bg/bg1.jpg";
import bg2 from "/public/images/bg/bg2.jpg";
import bg3 from "/public/images/bg/bg3.jpg";
import bg4 from "/public/images/bg/bg4.jpg";
import bg5 from "/public/images/bg/bg5.jpg";
import bg6 from "/public/images/bg/bg6.jpg";
import bg7 from "/public/images/bg/bg7.jpg";
import bg8 from "/public/images/bg/bg8.jpg";
import bg9 from "/public/images/bg/bg9.jpg";

export const systemSettingList = [
    {
        key: "system",
        title: "系统",
        desc: "存储、备份还原、用户角色",
        icon: "system",
        content: "SetSystem",
        children: [
            // {
            //     key: "user",
            //     title: "用户角色",
            //     icon: "users",
            //     content: "system/SetUser",
            // },
            // {
            //     key: "storage",
            //     title: "存储配置",
            //     icon: "disknet",
            //     content: "system/SetStorage",
            // },
            {
                key: "editor",
                title: "编辑器类型",
                icon: "note",
                content: "system/SetEditor",
            },
            {
                key: "backup",
                title: "备份还原",
                icon: "restore",
                content: "system/SetBackup",
            },
            {
                key: "password",
                title: "文件密码箱",
                icon: "pwdbox",
                content: "system/SetPassword",
            },
        ],
    },
    {
        key: "custom",
        title: "代理",
        desc: "本地代理、远程代理",
        icon: "proxy",
        content: "SetCustom",
        children: [
            {
                key: "local",
                title: "本地代理",
                icon: "local",
                content: "nas/SetLocal",
            },
            {
                key: "remote",
                title: "远程代理",
                icon: "netproxy",
                content: "nas/SetRemote",
            }

        ],
    },
    {
        key: "nas",
        title: "存储",
        desc: "NAS/webdav服务",
        icon: "nas",
        content: "SetNas",
        children: [
            {
                key: "localnas",
                title: "本地存储",
                icon: "disk",
                content: "nas/SetNas",
            },
            {
                key: "netnas",
                title: "远程存储",
                icon: "disknet",
                content: "nas/SetNas",
            }
        ],
    },
    {
        key: "account",
        title: "屏幕",
        desc: "壁纸/语言/锁屏/广告",
        icon: "zhuomian",
        content: "SetAccount",
        children: [
            {
                key: "wallpaper",
                title: "壁纸",
                icon: "style",
                content: "account/SetWallpaper",
            },
            {
                key: "language",
                title: "语言",
                icon: "lang",
                content: "account/SetLanguage",
            },
            {
                key: "lock",
                title: "锁屏",
                icon: "lock",
                content: "account/SetLock",
            },
            // {
            //     key: "ad",
            //     title: "广告",
            //     icon: "ad",
            //     content: "account/SetAd",
            // },
        ],
    },
];
export const settingsConfig = {
    background: {
        type: "image",
        color: "#ffffff",
        imageList: [bg1, bg2, bg3, bg4, bg5, bg6, bg7, bg8, bg9],
        url: bg6,
    },
    lock: {
        timeout: 0,
        activeTime: 0,
        password: ''
    },
    system: {
        userType: "person",
        storeType:  "local",
        netUrl: "",
    },
    wallpaper: {
        type: "image",
        color: "#ffffff",
        imageList: [bg1, bg2, bg3, bg4, bg5, bg6, bg7, bg8, bg9],
        url: bg6,
    },
    storeType: "local",
    storePath: "",
    netPort: "",
    netPath: "",
    storenet: {
        url: "",
        isCors: false,
    },
    webdavClient: {
        url: "",
        username: "",
        password: "",
    },
    editorType: "onlyoffice",
    onlyoffice: {
        url: "",
    },
    userInfo: {
        url: "",
        username: "",
        password: "",
    },
    proxyData: {
        proxyType: "http",
        domain: "",
        path: "",
    },
};