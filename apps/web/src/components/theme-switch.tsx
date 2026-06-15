import { useEffect } from 'react'
import { Check, Moon, Sun, Palette } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useTheme } from '@/context/theme-provider'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const colorThemes = [
  { value: 'violet', label: 'Violet', color: 'bg-violet-500' },
  { value: 'red', label: 'Red', color: 'bg-red-500' },
  { value: 'rose', label: 'Rose', color: 'bg-rose-500' },
  { value: 'sky', label: 'Sky', color: 'bg-sky-500' },
  { value: 'teal', label: 'Teal', color: 'bg-teal-500' },
  { value: 'yellow', label: 'Yellow', color: 'bg-yellow-500' },
  { value: 'amber', label: 'Amber', color: 'bg-amber-500' },
  { value: 'blue', label: 'Blue', color: 'bg-blue-500' },
  { value: 'cyan', label: 'Cyan', color: 'bg-cyan-500' },
  { value: 'emerald', label: 'Emerald', color: 'bg-emerald-500' },
  { value: 'fuchsia', label: 'Fuchsia', color: 'bg-fuchsia-500' },
  { value: 'green', label: 'Green', color: 'bg-green-500' },
  { value: 'indigo', label: 'Indigo', color: 'bg-indigo-500' },
  { value: 'lime', label: 'Lime', color: 'bg-lime-500' },
  { value: 'orange', label: 'Orange', color: 'bg-orange-500' },
  { value: 'pink', label: 'Pink', color: 'bg-pink-500' },
  { value: 'purple', label: 'Purple', color: 'bg-purple-500' },
] as const

export function ThemeSwitch() {
  const { theme, setTheme } = useTheme()

  /* Update theme-color meta tag
   * when theme is updated */
  useEffect(() => {
    let themeColor = '#fff'
    if (theme === 'dark') themeColor = '#020817'
    const metaThemeColor = document.querySelector("meta[name='theme-color']")
    if (metaThemeColor) metaThemeColor.setAttribute('content', themeColor)
  }, [theme])

  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger asChild>
        <Button variant='ghost' size='icon' className='scale-95 rounded-full'>
          <Sun className='size-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90' />
          <Moon className='absolute size-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0' />
          <span className='sr-only'>Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='end'>
        <DropdownMenuItem onClick={() => setTheme('light')}>
          <Sun className='me-2 size-4' />
          Light{' '}
          <Check
            size={14}
            className={cn('ms-auto', theme !== 'light' && 'hidden')}
          />
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('dark')}>
          <Moon className='me-2 size-4' />
          Dark
          <Check
            size={14}
            className={cn('ms-auto', theme !== 'dark' && 'hidden')}
          />
        </DropdownMenuItem>
        <DropdownMenuSub>
          <DropdownMenuSubTrigger>
            <Palette className='me-2 size-4' />
            <span>Colors</span>
          </DropdownMenuSubTrigger>
          <DropdownMenuSubContent>
            {colorThemes.map(({ value, label, color }) => (
              <DropdownMenuItem
                key={value}
                onClick={() => setTheme(value)}
              >
                <span className={cn('me-2 size-4 rounded-full', color)} />
                {label}
                <Check
                  size={14}
                  className={cn('ms-auto', theme !== value && 'hidden')}
                />
              </DropdownMenuItem>
            ))}
          </DropdownMenuSubContent>
        </DropdownMenuSub>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => setTheme('system')}>
          System
          <Check
            size={14}
            className={cn('ms-auto', theme !== 'system' && 'hidden')}
          />
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
