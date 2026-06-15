import { Globe } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { SUPPORTED_LANGUAGES } from '@/i18n/types'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

export function LanguageSwitcher() {
  const { i18n } = useTranslation()

  const handleLanguageChange = async (newLanguage: string) => {
    // eslint-disable-next-line no-console
    console.log('LanguageSwitcher: Changing language to', newLanguage, 'current:', i18n.language)
    await i18n.changeLanguage(newLanguage)
    // eslint-disable-next-line no-console
    console.log('LanguageSwitcher: Changed to', i18n.language)
  }

  return (
    <div className='flex items-center gap-2'>
      <Globe className='h-4 w-4 text-muted-foreground' />
      <Select value={i18n.language} onValueChange={handleLanguageChange}>
        <SelectTrigger className='w-[140px]'>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {Object.entries(SUPPORTED_LANGUAGES).map(([code, { name }]) => (
            <SelectItem key={code} value={code}>
              {name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}
