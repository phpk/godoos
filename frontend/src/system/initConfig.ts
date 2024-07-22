import { SystemOptions } from './type/type';
import { getSystemConfig } from './config'
const systemConfig = getSystemConfig()
export const defaultConfig: SystemOptions = {
  background: '#3A98CE',
  lang: systemConfig.lang,
  builtinFeature: [
    'MyComputer',
    'AppStore',
    'DataTimeTray',
    'BatteryTray',
    'NetworkTray',
    'ScreenshortTray',
    'ScreenRecorderTray',
    'ImageOpener',
    'UrlOpener',
    'TextOpener',
    'ShortCutOpener',
    'ExeOpener',
  ],
  userLocation: '/C/Users/',
  systemLocation: '/C/System/',
  login: {
    username: systemConfig.account.username,
    password: systemConfig.account.password,
    init: () => {
      return systemConfig.account.username === '';
    },
  },
  async loginCallback(username : string, password : string) {
    return (
      username === systemConfig.account.username &&
      password === systemConfig.account.password
    );
  },
};
