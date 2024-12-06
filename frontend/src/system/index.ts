import { markRaw, nextTick } from 'vue';
import { version } from '../../package.json';
import { useOsFile } from './core/FileOs';
import { OsFileSystem } from './core/FileSystem';
import { Eventer, initEventer, initEventListener } from './event';
import { initRootState, RootState } from './root';
import { SystemStateEnum } from './type/enum';
import {
  Saveablekey,
  Setting,
  SystemOptions,
  SystemOptionsCertainly,
  WinAppOptions,
} from './type/type';
import { BrowserWindow, BrowserWindowOption } from './window/BrowserWindow';

import { useAssistantStore } from '@/stores/assistant';
import { useModelStore } from "@/stores/model.ts";
import { useUpgradeStore } from '@/stores/upgrade';
import { RestartApp } from '@/util/goutil';
import { notifyError } from '@/util/msg';
import { isShareFile } from '@/util/sharePath';
import { pick } from '../util/modash';
import { clearSystemConfig, fetchGet, getClientId, getFileUrl, getSystemConfig, getSystemKey, setSystemConfig, setSystemKey } from './config';
import { OsFileInterface } from './core/FIleInterface';
import { extname } from './core/Path';
import { initBuiltinApp, initBuiltinFileOpener } from './initBuiltin';
import { defaultConfig } from './initConfig';
import { Tray, TrayOptions } from './menu/Tary';
import { Notify, NotifyConstructorOptions } from './notification/Notification';
import { Dialog } from './window/Dialog';

export type OsPlugin = (system: System) => void;
export type FileOpener = {
  name?: string;
  icon: string;
  hiddenInChosen?: boolean;
  func: (path: string, content: string) => void;
};
export class Bios {
  public static _onOpen: ((system: System) => void) | null = null;
  public static onOpen(func: (system: System) => void) {
    this._onOpen = func;
  }
  constructor() {
    //
  }
}

/**
 * @description: System 类，在初始化的过程中需要提供挂载点，以及一些配置
 */
export class System {
  public static GLOBAL_SYSTEM: System;

  readonly _options: SystemOptions;

  _rootState: RootState;
  private _eventer: Eventer;
  private _ready: ((value: System) => void) | null = null;
  private _error: ((reason: unknown) => void) | null = null;
  private _flieOpenerMap: Map<string, FileOpener> = new Map();
  version = version;
  isFirstRun = true;
  rootRef: HTMLElement | undefined = undefined;
  fs!: any;
  isReadyUpdateAppList = false;

  constructor(options?: SystemOptions) {

    this._options = this.initOptions(options);
    this._rootState = this.initRootState();
    System.GLOBAL_SYSTEM = this; // 挂载全局系统
    Bios._onOpen && Bios._onOpen(this);
    this._eventer = this.initEvent();
    this.firstRun();
    this.initSystem();
  }

  /**
   * @description: pure 初始化配置选项
   */
  private initOptions(options?: SystemOptions) {
    const tempOptions = Object.assign({}, defaultConfig, options);
    return tempOptions;
  }
  /**
   * @description: 获取系统配置
   */
  private initRootState(): RootState {
    return initRootState(this._options);
  }
  /**
   * @description: 初始化系统
   */
  private async initSystem() {
    this._rootState.state = SystemStateEnum.opening;
    this.fs = useOsFile(); // 初始化文件系统
    initBuiltinApp(this);
    initBuiltinFileOpener(this); // 注册内建文件打开器
    await this.initSavedConfig(); // 初始化保存的配置
    // 判断是否登录
    await this.isLogin();
    initEventListener(); // 初始化事件侦听


    this.initBackground(); // 初始化壁纸
    this.refershAppList();
    this.checkMessages();

    this.emit('start');
    this._ready && this._ready(this);

  }
  private checkMessages() {
    const config = getSystemConfig();
    const upgradeStore = useUpgradeStore();
    setTimeout(() => {
      //console.log(config.userType)
      if (config.userType == 'person') {
        upgradeStore.systemMessage();
      }else{
        upgradeStore.onlineMessage()
      }    
    }, 3000);
    

  }
  /**
   * @description: 判断是否登录
   */
  private async isLogin() {
    const config = getSystemConfig();

    if (config.userType == 'person') {
      if (!this._options.login) {
        this._rootState.state = SystemStateEnum.open;
        return;
      } else {
        if (this._options.login.init?.()) {
          this._rootState.state = SystemStateEnum.open;
          return;
        }

        this._rootState.state = SystemStateEnum.lock;
        const tempCallBack = this._options.loginCallback;
        if (!tempCallBack) {
          throw new Error('没有设置登录回调函数');
        }
        this._options.loginCallback = async (username: string, password: string) => {
          const res = await tempCallBack(username, password);
          if (res) {
            this._rootState.state = SystemStateEnum.open;
            return true;
          }
          return false;
        };
      }
    } else {
      // 非个人用户
      const userInfo = config.userInfo;
      if (userInfo.url == '') {
        return true;
      }
      const res = await fetchGet(`${userInfo.url}/member/islogin`);
      if (!res.ok) {
        notifyError('登录失败，请检查网络连接或重新登录');
        return true;
      }
      const data = await res.json();
      //console.log(data)
      if (data.success) {
        // 登录成功
        this._rootState.state = SystemStateEnum.open;
        return true;
      } else {
        // 登录失败
        //return true;
        this._rootState.state = SystemStateEnum.lock;
        // 登录回调
        this._options.loginCallback = async (
          username: string, password: string, loginCode?: {
            github_code?: string,
            gitee_code?: string
          }
        ) => {
          const serverUrl = config.userInfo.url + '/member/login'
          const res: any = await fetch(serverUrl, {
            method: "POST",
            body: JSON.stringify({
              username: username,
              password: password,
              github_code: loginCode?.github_code,
              gitee_code: loginCode?.gitee_code,
              clientId: getClientId(),
            }),
          });
          if (!res.ok) {
            return false
          }

          const jsondata = await res.json();
          if (jsondata.success) {
            jsondata.data.url = config.userInfo.url
            jsondata.data.password = password
            config.userInfo = jsondata.data
            setSystemConfig(config);
            this._rootState.state = SystemStateEnum.open;
            this.refershAppList()
            window.location.href = "/";
            return true
          } else {
            notifyError(jsondata.message)
            return false
          }
        };



        // //const tmpCallBack = this._options.authCallback;
        // this._options.authCallback = async (username: string, password: string, captcha :string) => {
        //   //this._rootState.state = SystemStateEnum.open;
        //     return false;
        // };
        // return false;
      }
    }
  }

  /**
   * @description: 初始化壁纸
   */
  initBackground() {
    const background = getSystemKey('background');
    if (background.type === 'image') {
      this._rootState.options.background = background.url
    } else {
      this._rootState.options.background = background.color;
    }
  }
  /**
   * @description: 初始化事件系统
   */
  private initEvent() {
    return initEventer();
  }

  refershAppList() {

    const system = useSystem();
    if (!system) return;
    const fileUrl = getFileUrl();
    if (!fileUrl) return;
    try {
      fetchGet(`${fileUrl}/desktop`).then(res => res.json()).then(res => {
        if (res && res.code == 0) {
          system._rootState.apps.splice(0, system._rootState.apps.length, ...res.data.apps);
          system._rootState.menulist.splice(0, system._rootState.menulist.length, ...res.data.menulist);
        }
      })
    } catch (error) {
      console.log(error)
    }
    
  }


  initAppList() {
    this.isReadyUpdateAppList = true;
    nextTick(() => {
      if (this.isReadyUpdateAppList) {
        this.isReadyUpdateAppList = false;
        this.refershAppList();
      }
    });
  }

  replaceFileSystem(fs: OsFileInterface) {
    this.fs = fs;
    //this.initAppList();
    this.refershAppList();
  }
  mountVolume(path: string, fs: OsFileInterface) {
    if (this.fs instanceof OsFileSystem) {
      this.fs.mountVolume(path, fs);
    } else {
      console.error('自定义文件系统不支持挂载卷');
    }
  }

  /**
   * @description: 初始化保存的配置
   */
  private async initSavedConfig() {
    //const config = await this.fs.readFile(join(this._options.systemLocation || '', 'os/config.json'));
    const config = getSystemConfig();
    if (config) {
      try {
        this._rootState.options = Object.assign(this._rootState.options, config);
      } catch {
        new Notify({
          title: '提示',
          content: '配置文件格式错误',
        });
      }
    }
  }
  setConfig<T extends keyof SystemOptionsCertainly>(key: T, value: SystemOptionsCertainly[T]): Promise<void>;
  setConfig<T extends string>(
    key: T,
    value: T extends keyof SystemOptionsCertainly ? SystemOptionsCertainly[T] : unknown
  ): Promise<void>;
  setConfig<T extends keyof SystemOptionsCertainly>(key: string, value: SystemOptionsCertainly[T]) {
    this._rootState.options[key] = value;
    if (Saveablekey.includes(key as any)) {
      return setSystemConfig(pick(this._rootState.options, ...Saveablekey))
    } else {
      return Promise.resolve();
    }
  }

  getConfig<T extends keyof SystemOptionsCertainly>(key: T): SystemOptionsCertainly[T];
  getConfig<T extends string>(key: T): unknown;
  getConfig(key: string) {
    return this._rootState.options[key];
  }

  private addWindowSysLink(loc: string, options: any, force = false) {
    if (force) {
      this.fs.writeFile(
        `${this._options.userLocation}${loc}/` + options.name + '.exe',
        `link::${loc}::${options.name}::${options.icon}`
      );
    }
    //if(!options.window.content)return
    if (typeof options.window.content !== 'string') {
      console.log('index.ts:', options.window.content)
      //options.window.content = markRaw(options.window.content);
    }
    // console.log(options.window)
    // console.log(options.name)
    this._rootState.windowMap[loc].set(options.name, options);
  }

  /**
   * @description: 添加应用
   * force 表示强制，在每次启动时都会添加
   */
  addApp(options: WinAppOptions, force = false) {
    this.addWindowSysLink('Desktop', options, force);
  }
  addMenuList(options: WinAppOptions, force = false) {
    this.addWindowSysLink('Menulist', options, force);
  }
  addBuiltInApp(options: WinAppOptions) {
    this._rootState.windowMap['Builtin'].set(options.name, options);
  }

  whenReady(): Promise<System> {
    return new Promise<System>((resolve, reject) => {
      this._ready = resolve;
      this._error = reject;
    });
  }
  firstRun() {
    if (getSystemKey('isFirstRun')) {
      this.isFirstRun = false;
      return false;
    } else {
      this.isFirstRun = true;
      setSystemKey('isFirstRun', true)
      this.emit('firstRun');
      const promptStore = useAssistantStore()
      promptStore.initPrompt()
      const modelStore = useModelStore()
      modelStore.initModel()
      return true;
    }
  }
  shutdown() {
    this._rootState.state = SystemStateEnum.close;
  }
  reboot() {
    this._rootState.state = SystemStateEnum.close;
    RestartApp();
  }
  recover() {
    clearSystemConfig()
    this._rootState.state = SystemStateEnum.close;
    this.fs.removeFileSystem().then(() => {
      window.indexedDB.deleteDatabase("GodoDatabase");

      RestartApp();
    })

  }
  getEventer() {
    return this._eventer;
  }
  emit(event: string, ...args: any[]) {
    this.emitEvent(event, ...args);
  }
  emitEvent(event: string, ...args: any[]) {
    const eventArray = event.split('.');
    eventArray.forEach((_: any, index) => {
      const tempEvent = eventArray.slice(0, index + 1).join('.');
      this._eventer.emit(tempEvent, event, args);
    });
    this._eventer.emit('system', event, args);
  }
  on(event: string, callback: (...args: any[]) => void): void {
    this.mountEvent(event, callback);
  }
  mountEvent(event: string | string[], callback: (...args: any[]) => void) {
    if (Array.isArray(event)) {
      event.forEach((item) => {
        this.mountEvent(item, callback);
      });
      return;
    } else {
      this._eventer.on(event, callback);
    }
  }

  offEvent(event?: string, callback?: (...args: any[]) => void): void {
    this._eventer.off(event, callback);
  }

  /** 注册文件打开器 */
  registerFileOpener(type: string | string[], opener: FileOpener) {
    if (Array.isArray(type)) {
      type.forEach((item) => {
        this._flieOpenerMap.set(item, opener);
      });
      return;
    }
    this._flieOpenerMap.set(type, opener);
  }
  getOpener(type: string) {
    return this._flieOpenerMap.get(type);
  }
  getAllFileOpener() {
    return this._flieOpenerMap;
  }

  /** 注册设置app的设置页面 */
  registerSettingPanel(setting: Setting) {
    const temp = {
      ...setting,
      content: markRaw(setting.content),
    };
    this._rootState.settings?.push(temp);
  }
  /**打开os 文件系统的文件 */
  async openFile(path: string) {
      const fileStat = await this.fs.stat(path)
      //console.log(fileStat)
      if (!fileStat) {
        throw new Error('文件不存在');
      }
      // 如果fileStat为目录
      if (fileStat?.isDirectory) {
        // 从_fileOpenerMap中获取'link'对应的函数并调用
        this._flieOpenerMap.get('dir')?.func.call(this, path, '');
        return;
      } else {
        const header = {
          pwd: ''
        }
        // 读取文件内容
        let fileContent = await this.fs.readFile(path, header);
        // 改文件需要输入密码
        if (fileContent && fileContent.code == -1 && fileContent.error == 'needPwd') {
          const temp = await Dialog.showInputBox()
          if (temp.response !== 1) {
            return
          }
          header.pwd = temp?.inputPwd ? temp?.inputPwd : ''
          const reOpen = await this.fs.readFile(path, header);
          //console.log(reOpen)
          if (reOpen === false || reOpen.code === -1) {
            notifyError("文件密码错误")
            return
          }
          fileContent = reOpen
        }
        //console.log(fileStat)
        // 从_fileOpenerMap中获取文件扩展名对应的函数并调用
        const fileName = extname(fileStat?.name || '') || 'link'
        //console.log(fileName,fileContent)
        this._flieOpenerMap
          .get(fileName)
          ?.func.call(this, path, fileContent || '');
    }
  }
  // 插件系统
  use(func: OsPlugin): void {
    return func(this);
  }
  // 状态序列化和反序列化
  async serializeState(): Promise<string> {
    const serializeFile = await this.fs.serializeFileSystem();
    return JSON.stringify(serializeFile);
  }
  deserializeState(state: string) {
    this.fs.deserializeFileSystem(JSON.parse(state));
  }

  outerFileDropCallback:
    | ((path: string, list: FileList | undefined, process: (path: string) => void) => void)
    | null = null;
  // 当从外部拖入文件时
  onOuterFileDrop(func: (path: string, list: FileList | undefined, process: (path: string) => void) => void) {
    this.outerFileDropCallback = func;
  }
  /** 方便的通过system创建window */
  createWindow(options: BrowserWindowOption) {
    const win = new BrowserWindow(options);
    return win;
  }
  /** 方便的通过system创建notify */
  createNotify(options: NotifyConstructorOptions) {
    return new Notify(options);
  }
  /** 方便的通过system创建Dialog */
  createDialog() {
    return Dialog;
  }
  /** 方便的通过system创建Tray */
  createTray(options: TrayOptions) {
    return new Tray(options);
  }

  errorHandler = 0;
  emitError(error: string) {
    this._error && this._error(error);
    this._rootState.error = error;
    this.errorHandler = Date.now();
    setTimeout(() => {
      if (Date.now() - this.errorHandler > 1000 * 3) {
        this._rootState.error = '';
      }
    }, 1000 * 4);
  }
}
export function useSystem() {
  return System.GLOBAL_SYSTEM!;
}

export { t } from '../i18n';
export { dealIcon } from '../util/Icon';
export type { OsFileInterface } from './core/FIleInterface';
export * from './core/FileSystem';
export * from './core/Path';
export { Menu } from './menu/Menu';
export { MenuItem } from './menu/MenuItem';
export { Tray } from './menu/Tary';
export { Notify } from './notification/Notification';
export type { SystemOptions, WinApp } from './type/type';
export { BrowserWindow } from './window/BrowserWindow';
export { Dialog } from './window/Dialog';
export { vDragable } from './window/MakeDragable';

