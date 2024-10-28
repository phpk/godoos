import { createI18n } from 'vue-i18n'
import elEnLocale from 'element-plus/es/locale/lang/en'
import elZhLocale from 'element-plus/es/locale/lang/zh-cn'
import zhLang  from './lang/zh.json';
import enLang  from './lang/en.json';
import { getSystemKey, setSystemKey } from '@/system/config'
export function getLang() {
  let currentLang = getSystemKey('lang')
  if (!currentLang) {
    try {
      const supported = ["en", "zh-cn"]
      const browserLang = (navigator.language || (navigator as any).browserLanguage).toLowerCase()
      if (supported.includes(browserLang)) {
        currentLang = browserLang
      } else {
        currentLang = "en";
      }
      setLang(currentLang)
    } catch (e) {
      console.log(e);
    }
  }
  return currentLang
}


const messages = {
  en: {
    ...enLang,
    ...elEnLocale
  },
  'zh-cn': {
    ...zhLang,
    ...elZhLocale
  }
}
export const i18n = createI18n({
  globalInjection: true,
  legacy: false,
  locale: getLang(),
  messages: messages
})
export function setLang(lang: string) {
  //currentLang = lang
  setSystemKey('lang', lang)
}
export function changeLang() {
  const lang = getLang()
  const setlang = lang == 'en' ? 'zh-cn' : 'en'
  setLang(setlang)
  return setlang
}
export function t(textkey :string) {
  return i18n.global.t(textkey)
}
export function dealSystemName(name: string) {
  const sysNames:any = {
    "C": t('system'),
    "D": t('document'),
    "E": t('office'),
    "B": t('recycle'),
    "F": t('share'),
    "myshare": t('myshare'),
    "othershare": t('othershare')
  }
  if (sysNames[name]) {
    return sysNames[name];
  } else {
    return name;
  }
  //return name;
}