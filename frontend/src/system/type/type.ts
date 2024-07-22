import { BrowserWindowOption } from '../window/BrowserWindow';
import { MenuItem, MenuItemConstructorOptions } from '../menu/MenuItem';
export type BuiltinFeature =
  | 'MyComputer'
  | 'AppStore'
  | 'DataTimeTray'
  | 'BatteryTray'
  | 'NetworkTray'
  | 'ScreenshortTray'
  | 'ScreenRecorderTray'
  | 'ImageOpener'
  | 'UrlOpener'
  | 'TextOpener'
  | 'ShortCutOpener'
  | 'ExeOpener';
export interface InitFileItem {
  type: string;
  name: string;
  children?: InitFileItem[];
  mode?: number;
  content?: string;
}
export interface Setting {
  key: string;
  title: string;
  desc: string;
  icon: string;
  content: any;
}
export type SafeAny = unknown;
export const Saveablekey: ('lang' | 'logo' | 'background' | 'rootStyle')[] = [
  'lang',
  'logo',
  'background',
  'rootStyle',
];
export type SystemOptionsSaveable = Pick<SystemOptionsCertainly, (typeof Saveablekey)[number]>;

export interface SystemOptionsCertainly {
  lang?: string;
  logo?: string;
  background?: string;
  rootStyle?: any;
  builtinFeature?: BuiltinFeature[];
  desktop?: WinAppOptions[];
  magnet?: WinAppOptions[];
  menulist?: WinAppOptions[];
  fs?: any;
  userLocation?: string;
  systemLocation?: string;
  initFile?: InitFileItem;
  login?: {
    username: string;
    password: string;
    init?: () => boolean;
  };
  contextMenus?: Array<MenuItemConstructorOptions | MenuItem>;
  noPassword?: boolean;
  loginCallback?: (username: string, password: string) => Promise<boolean>;
}
export type SystemOptions = SystemOptionsCertainly & {
  [key: string]: SafeAny;
};
export interface WinApp {
  icon: string;
  path: string;
}
export interface WinAppOptions {
  name: string;
  icon?: string;
  multiple?: boolean;
  window: BrowserWindowOption;
  _hasShow?: boolean;
}
