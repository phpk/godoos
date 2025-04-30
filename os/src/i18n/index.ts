// src/i18n/index.ts
import { createI18n } from 'vue-i18n';
import elEnLocale from 'element-plus/es/locale/lang/en';
import elZhLocale from 'element-plus/es/locale/lang/zh-cn';
import zhLang from './lang/zh.json';
import enLang from './lang/en.json';

export const supported = ["en", "zh-cn"];
export const languages = [
  { value: "en", label: "English" },
  { value: "zh-cn", label: "中文" },
];

export function getLang() {
  let currentLang:any = localStorage.getItem('godoos_lang');
  if (!currentLang) {
    try {
      const browserLang = (navigator.language || (navigator as any).browserLanguage).toLowerCase();
      if (supported.includes(browserLang)) {
        currentLang = browserLang;
      } else {
        currentLang = "zh-cn";
      }
      setLang(currentLang);
    } catch (e) {
      console.log(e);
    }
  }
  return currentLang || 'zh-cn';
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
};

export const i18n = createI18n({
  globalInjection: true,
  legacy: false,
  locale: getLang(),
  messages: messages,
  silentTranslationWarn: true,
  missing: (_, key) => {
    return key; // 返回键本身或默认值
  },
});

export function setLang(lang: any) {
  if (supported.includes(lang)) {
    localStorage.setItem('godoos_lang', lang);
    i18n.global.locale.value = lang;
  }
}

export function changeLang() {
  const lang = getLang();
  const setlang = lang === 'en' ? 'zh-cn' : 'en';
  setLang(setlang);
  return setlang;
}

export function t(textkey: string) {
  return i18n.global.t(textkey);
}

export function dealSystemName(name: string) {
  const sysNames: any = {
    "C": t('system') + '盘',
    "D": t('document') + '盘',
    "E": t('office') + '盘',
    "B": t('recycle'),
    "F": t('share') + '盘',
    "myshare": t('myshare'),
    "othershare": t('othershare')
  };
  return sysNames[name] || t(name);
}