import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import {
  SUPPORTED_LANGUAGES,
  DEFAULT_LANGUAGE,
  type SupportedLanguage,
} from './types'

// 导入翻译文件（静态导入用于基础翻译）
import enCommon from './locales/en/common.json'
import enAuth from './locales/en/auth.json'
import enErrors from './locales/en/errors.json'
import zhCNCommon from './locales/zh-CN/common.json'
import zhCNAuth from './locales/zh-CN/auth.json'
import zhCNErrors from './locales/zh-CN/errors.json'

// 从 localStorage 读取保存的语言
const getSavedLanguage = (): SupportedLanguage => {
  try {
    const saved = localStorage.getItem('i18next_lng')
    if (saved && saved in SUPPORTED_LANGUAGES) {
      return saved as SupportedLanguage
    }
  } catch {
    // Ignore localStorage errors
  }
  return DEFAULT_LANGUAGE
}

// 资源配置
const resources = {
  en: {
    common: enCommon,
    auth: enAuth,
    errors: enErrors,
  },
  'zh-CN': {
    common: zhCNCommon,
    auth: zhCNAuth,
    errors: zhCNErrors,
  },
}

// 使用默认的 i18next 实例（不创建新实例）
i18n.use(initReactI18next).init({
  // 资源
  resources,

  // 从 localStorage 读取语言或使用默认语言
  lng: getSavedLanguage(),
  fallbackLng: DEFAULT_LANGUAGE,

  // 命名空间
  ns: ['common', 'auth', 'errors'],
  defaultNS: 'common',

  // React 配置
  react: {
    useSuspense: false,
    bindI18n: 'languageChanged',
    bindI18nStore: '',
    transEmptyNodeValue: '',
    transSupportBasicHtmlNodes: true,
    transKeepBasicHtmlNodesFor: ['br', 'strong', 'i', 'em'],
  },

  // 调试
  debug: import.meta.env.DEV,

  // 插值
  interpolation: {
    escapeValue: false,
  },
})

// eslint-disable-next-line no-console
console.log('i18n initialized with language:', i18n.language)

// 语言变化时更新文档方向和保存到 localStorage
i18n.on('languageChanged', (lng) => {
  // eslint-disable-next-line no-console
  console.log('languageChanged to:', lng, 'i18n.language:', i18n.language)
  const lang = lng as SupportedLanguage
  const dir = SUPPORTED_LANGUAGES[lang]?.dir || 'ltr'
  document.documentElement.setAttribute('dir', dir)
  document.documentElement.setAttribute('lang', lng)
  // 保存到 localStorage
  try {
    localStorage.setItem('i18next_lng', lng)
  } catch {
    // Ignore localStorage errors
  }
})

export default i18n
