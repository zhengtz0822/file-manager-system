import { create } from 'zustand'
import i18n from '@/i18n/config'
import {
  SUPPORTED_LANGUAGES,
  DEFAULT_LANGUAGE,
  type SupportedLanguage,
} from '@/i18n/types'

interface I18nState {
  language: SupportedLanguage
  setLanguage: (language: SupportedLanguage) => Promise<void>
  isRTL: boolean
}

export const useI18nStore = create<I18nState>((set) => {
  // 初始化当前语言
  const initLanguage =
    (i18n.language as SupportedLanguage) || DEFAULT_LANGUAGE

  // 监听 i18next 语言变化，同步到 Zustand store
  i18n.on('languageChanged', (lng) => {
    const lang = lng as SupportedLanguage
    set({
      language: lang,
      isRTL: SUPPORTED_LANGUAGES[lang]?.dir === 'rtl',
    })
  })

  return {
    language: initLanguage,
    isRTL: SUPPORTED_LANGUAGES[initLanguage]?.dir === 'rtl',

    setLanguage: async (language: SupportedLanguage) => {
      await i18n.changeLanguage(language)
      // i18next 的 languageChanged 事件会自动更新 store
    },
  }
})
