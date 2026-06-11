import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import enUS from './locales/en-US'

export const LOCALES = {
  zhCN: 'zh-CN',
  enUS: 'en-US',
}

export const getStoredLocale = () => localStorage.getItem('app-locale') || LOCALES.zhCN

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: getStoredLocale(),
  fallbackLocale: LOCALES.zhCN,
  messages: {
    [LOCALES.zhCN]: zhCN,
    [LOCALES.enUS]: enUS,
  },
})

export function setI18nLocale(locale) {
  i18n.global.locale.value = locale
  localStorage.setItem('app-locale', locale)
  document.documentElement.setAttribute('lang', locale)
}

export function t(key, params) {
  return i18n.global.t(key, params)
}

setI18nLocale(getStoredLocale())
