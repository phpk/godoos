import { zhCN } from './system/zh';
import { enUS } from './system/en';
import { getSystemKey, setSystemKey } from '../system/config'

function getLang() {
  let currentLang = getSystemKey('lang')
  if (!currentLang) {
    try {
      const supported = ["en", "zh"]
      const { 0: browserLang } = navigator.language.split("-");
      if (supported.includes(browserLang)) {
        currentLang = browserLang
      } else {
        currentLang = "en";
      }
    } catch (e) {
      console.log(e);
    }
  }
  return currentLang
}

export let currentLang = getLang()
export function setLang(lang: string) {
  currentLang = lang
  setSystemKey('lang', lang)
}

export function i18n(key: string) {
  if (currentLang == 'zh') {
    return zhCN[key] ?? key;
  } else if (currentLang == 'en') {
    return enUS[key] ?? key;
  } else {
    return enUS[key] ?? key;
  }
}

export function dealSystemName(name: string) {
  const sysNames:any = {
    "C": i18n('system'),
    "D": i18n('document'),
    "E": i18n('office'),
    "B": i18n('recycle')
  }
  if (sysNames[name]) {
    return sysNames[name];
  } else {
    return name;
  }
  //return name;
}