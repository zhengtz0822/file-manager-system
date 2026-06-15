import { type ReactNode } from 'react'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n/config'

interface I18nProviderProps {
  children: ReactNode
}

export function I18nProvider({ children }: I18nProviderProps) {
  // eslint-disable-next-line no-console
  console.log('I18nProvider rendering with i18n instance:', i18n.language, 'instance:', !!i18n)
  return <I18nextProvider i18n={i18n}>{children}</I18nextProvider>
}
