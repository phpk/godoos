import { initRootState, RootState } from './root';
import { SystemStateEnum } from './type/enum';
import { markRaw, nextTick } from 'vue';
import {
  Saveablekey,
  Setting,
  SystemOptions,
  SystemOptionsCertainly,
  WinAppOptions,
} from './type/type';
import { initEventer, Eventer, initEventListener } from './event';
import { OsFileSystem } from './core/FileSystem';
import {useOsFile} from './core/FileOs';
import { version } from '../../package.json';
import { BrowserWindow, BrowserWindowOption } from './window/BrowserWindow';

import { extname } from './core/Path';
import { initBuiltinApp, initBuiltinFileOpener } from './initBuiltin';
import { defaultConfig } from './initConfig';
import { OsFileInterface } from './core/FIleInterface';
import { Notify, NotifyConstructorOptions } from './notification/Notification';
import { Dialog } from './window/Dialog';
import { pick } from '../util/modash';
import { Tray, TrayOptions } from './menu/Tary';
import { InitSystemFile, InitUserFile } from './core/SystemFileConfig';
import { createInitFile } from './core/createInitFile';
import { getSystemConfig, getSystemKey, setSystemKey, setSystemConfig, clearSystemConfig } from './config'


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
    //await startApp();
    /**
     * 过程：激活屏幕，桥接事件。
     */
    this._rootState.state = SystemStateEnum.opening;
    initBuiltinFileOpener(this); // 注册内建文件打开器
    await this.initFileSystem(); // 初始化文件系统
    await this.initSavedConfig(); // 初始化保存的配置
    initBuiltinApp(this); // 初始化内建应用
    this.initApp(); // 初始化配置应用到app文件夹中
    this.initAppList(); // 刷新app文件夹，展示应用
    // 判断是否登录
    this.isLogin();
    initEventListener(); // 初始化事件侦听

    this._ready && this._ready(this);
    //this.runPlugin(this); // 运行fs中插件
    this.initBackground(); // 初始化壁纸
    this.emit('start');
    
    // setTimeout(() => {
    //   if (this._rootState.magnet?.length < 1) {
    //     this.refershAppList()
    //   }
    // }, 6000);

  }
  /**
   * @description: 判断是否登录
   */
  private isLogin() {
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
  }

  /**
   * @description: 初始化壁纸
   */
  initBackground() {
    const background = getSystemKey('background');
    if (background.type === 'image') {
      this._rootState.options.background = background.url
    }else{
      this._rootState.options.background = background.color;
    }
    // this._rootState.options.background = '/image/bg/bg5.jpg';
  }
  /**
   * @description: 初始化事件系统
   */
  private initEvent() {
    /**
     * 过程：监听事件，处理事件
     */
    return initEventer();
  }
  private initApp() {

    this._rootState.options.desktop?.forEach((item) => {
      this.addApp(item);
    });
    this._rootState.options.magnet?.forEach((item) => {
      this.addMagnet(item);
    });
    this._rootState.options.menulist?.forEach((item) => {
      this.addMenuList(item);
    });

  }

  refershAppList() {
    const APP_TYPE = ['apps', 'magnet', 'menulist'];
    const system = useSystem();
    for (let i = 0; i < APP_TYPE.length; i++) {
      const element = APP_TYPE[i];
      system?.fs
        .readdir(
          `${system._options.userLocation}${{
            apps: 'Desktop',
            magnet: 'Magnet',
            menulist: 'Menulist',
          }[element]
          }`
        )
        .then((res:any) => {
          if (res) {
            const list = res;
            const tempList:any = [];
            for (let j = 0; j < list.length; j++) {
              const item = list[j];

              tempList.push(item);
            }

            switch (element) {
              case 'apps':
                useSystem()._rootState.apps.splice(0, useSystem()._rootState.apps.length, ...tempList);
                break;
              case 'magnet':
                useSystem()._rootState.magnet.splice(0, useSystem()._rootState.magnet.length, ...tempList);
                break;
              case 'menulist':
                useSystem()._rootState.menulist.splice(
                  0,
                  useSystem()._rootState.menulist.length,
                  ...tempList
                );
                break;
              default:
                break;
            }
          }
        });
    }
  }

  isReadyUpdateAppList = false;
  initAppList() {
    this.isReadyUpdateAppList = true;
    nextTick(() => {
      if (this.isReadyUpdateAppList) {
        this.isReadyUpdateAppList = false;
        this.refershAppList();
      }
    });
  }

  private async initFileSystem() {
    // 如果传入了自定义fs，就使用传入的fs
    //if (this._options.fs) {
    if (getSystemKey("storeType") == 'local'){
      //this.fs = this._options.fs;
      this.fs = useOsFile();
      await this.initOutSystem()
    } else {
      this.fs = await new OsFileSystem().initFileSystem(this._options);
      (this.fs as OsFileSystem).on('error', (err: string) => {
        this.emitError(err);
      });
      this.fs.registerWatcher(new RegExp(`^${this._options.userLocation}`), () => {
        this.initAppList();
      });
    }
  }
  async initOutSystem() {
    if (!this.isFirstRun) {
      return;
    }
    await this.fs.mkdir('/C');
    await this.fs.mkdir('/D');
    await this.fs.mkdir('/E');
    await this.fs.mkdir('/B');
    await createInitFile(this.fs, InitUserFile, this._options.userLocation);
    await createInitFile(this.fs, InitSystemFile, this._options.systemLocation);

    // await this.fs.chmod('/C', OsFileMode.ReadWrite);
    // await this.fs.chmod('/D', OsFileMode.ReadWrite);
    // await this.fs.chmod('/E', OsFileMode.ReadWrite);
    // await this.fs.chmod('/B', OsFileMode.ReadWrite);

  }
  replaceFileSystem(fs: OsFileInterface) {
    this.fs = fs;
    this.initAppList();
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
      // return this.fs.writeFile(
      //   join(this._options.systemLocation || '', 'os/config.json'),

      //   JSON.stringify(pick(this._rootState.options, ...Saveablekey))
      // );
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

  private addWindowSysLink(loc: string, options: WinAppOptions, force = false) {
    if (this.isFirstRun || force) {
      this.fs.writeFile(
        `${this._options.userLocation}${loc}/` + options.name + '.exe',
        `link::${loc}::${options.name}::${options.icon}`
      );
    } else {
      this.initAppList();
    }
    if (typeof options.window.content === 'string') {
      // TODO: 当content是string的时候
    } else {
      options.window.content = markRaw(options.window.content);
    }
    this._rootState.windowMap[loc].set(options.name, options);
  }

  // async runPlugin(system: System) {
  //   const pluginsFile = await this.fs.readdir(`${this._options.systemLocation}plugs`);
  //   if (pluginsFile) {
  //     await Promise.all(
  //       pluginsFile.map(async (file) => {
  //         const fileContent = await this.fs.readFile(file.path);
  //         if (file.isFile) {
  //           const content = fileContent;
  //           if (content) {
  //             //new Shell(system, '/', 'root').exec('node ' + file.path);
  //           }
  //         }
  //       })
  //     );
  //   }
  // }
  /**
   * @description: 添加应用
   * force 表示强制，在每次启动时都会添加
   */
  addApp(options: WinAppOptions, force = false) {
    this.addWindowSysLink('Desktop', options, force);
  }
  addMagnet(options: WinAppOptions, force = false) {
    this.addWindowSysLink('Magnet', options, force);
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
      // localStorage.setItem('godoOS_firstRun', 'true');
      this.emit('firstRun');
      return true;
    }
  }
  shutdown() {
    this._rootState.state = SystemStateEnum.close;
  }
  reboot() {
    this._rootState.state = SystemStateEnum.close;
    window.location.reload();
  }
  recover() {
    clearSystemConfig()
    this._rootState.state = SystemStateEnum.close;
    this.fs.removeFileSystem().then(() => {
      window.indexedDB.deleteDatabase("GodoDatabase");
      
      window.location.reload();
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
    eventArray.forEach((_:any, index) => {
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
    const fileStat = await this.fs.stat(path);
    if (!fileStat) {
      throw new Error('文件不存在');
    }
    // 如果fileStat为目录
    if (fileStat?.isDirectory) {
      // 从_fileOpenerMap中获取'link'对应的函数并调用
      this._flieOpenerMap.get('dir')?.func.call(this, path, '');
      return;
    } else {
      // 读取文件内容
      const fileContent = await this.fs.readFile(path);
      // 从_fileOpenerMap中获取文件扩展名对应的函数并调用
      this._flieOpenerMap
        .get(extname(fileStat?.path || '') || 'link')
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

export * from './core/Path';
export * from './core/FileSystem';
export { BrowserWindow } from './window/BrowserWindow';
export { Notify } from './notification/Notification';
export { Dialog } from './window/Dialog';
export type { SystemOptions, WinApp } from './type/type';
export { vDragable } from './window/MakeDragable';
export type { OsFileInterface } from './core/FIleInterface';
export { i18n } from '../i18n';
export { Tray } from './menu/Tary';
export { dealIcon } from '../util/Icon';
export { Menu } from './menu/Menu';
export { MenuItem } from './menu/MenuItem';