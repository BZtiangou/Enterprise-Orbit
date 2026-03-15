import { createI18n } from 'vue-i18n'
import zh from './zh'
import en from './en'

type MessageSchema = typeof zh

const messages = {
  zh,
  en
} as const

const savedLocale = (localStorage.getItem('locale') || 'zh') as 'zh' | 'en'

const i18n = createI18n<[MessageSchema], 'zh' | 'en'>({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'en',
  messages,
  globalInjection: true
})

export default i18n

export const setLocale = (locale: 'zh' | 'en') => {
  (i18n.global.locale as unknown as { value: 'zh' | 'en' }).value = locale
  localStorage.setItem('locale', locale)
  document.querySelector('html')?.setAttribute('lang', locale)
}

export const getLocale = (): 'zh' | 'en' => {
  return (i18n.global.locale as unknown as { value: 'zh' | 'en' }).value
}
