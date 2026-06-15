import type { Direction } from '@/context/direction-provider'

// 支持的语言配置
export const SUPPORTED_LANGUAGES = {
  en: { name: 'English', dir: 'ltr' as Direction },
  'zh-CN': { name: '简体中文', dir: 'ltr' as Direction },
  'zh-TW': { name: '繁體中文', dir: 'ltr' as Direction },
} as const

export type SupportedLanguage = keyof typeof SUPPORTED_LANGUAGES

// 默认语言
export const DEFAULT_LANGUAGE: SupportedLanguage = 'en'

// Cookie 键名
export const LANGUAGE_COOKIE_KEY = 'i18next_lng'
