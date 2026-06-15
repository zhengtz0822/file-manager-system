import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { ChevronDownIcon } from '@radix-ui/react-icons'
import { zodResolver } from '@hookform/resolvers/zod'
import { fonts } from '@/config/fonts'
import { showSubmittedData } from '@/lib/show-submitted-data'
import { cn } from '@/lib/utils'
import { useFont } from '@/context/font-provider'
import { useTheme } from '@/context/theme-provider'
import { Button, buttonVariants } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const appearanceFormSchema = z.object({
  theme: z.enum([
    'light',
    'dark',
    'violet',
    'red',
    'rose',
    'sky',
    'teal',
    'yellow',
    'amber',
    'blue',
    'cyan',
    'emerald',
    'fuchsia',
    'green',
    'indigo',
    'lime',
    'orange',
    'pink',
    'purple',
  ]),
  font: z.enum(fonts),
})

type AppearanceFormValues = z.infer<typeof appearanceFormSchema>

export function AppearanceForm() {
  const { font, setFont } = useFont()
  const { theme, setTheme } = useTheme()

  // This can come from your database or API.
  const defaultValues: Partial<AppearanceFormValues> = {
    theme: theme as 'light' | 'dark' | 'violet' | 'red' | 'rose' | 'sky' | 'teal' | 'yellow' | 'amber' | 'blue' | 'cyan' | 'emerald' | 'fuchsia' | 'green' | 'indigo' | 'lime' | 'orange' | 'pink' | 'purple',
    font,
  }

  const form = useForm<AppearanceFormValues>({
    resolver: zodResolver(appearanceFormSchema),
    defaultValues,
  })

  function onSubmit(data: AppearanceFormValues) {
    if (data.font != font) setFont(data.font)
    if (data.theme != theme) setTheme(data.theme)

    showSubmittedData(data)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-8'>
        <FormField
          control={form.control}
          name='font'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Font</FormLabel>
              <div className='relative w-max'>
                <FormControl>
                  <select
                    className={cn(
                      buttonVariants({ variant: 'outline' }),
                      'w-50 appearance-none font-normal capitalize',
                      'dark:bg-background dark:hover:bg-background'
                    )}
                    {...field}
                  >
                    {fonts.map((font) => (
                      <option key={font} value={font}>
                        {font}
                      </option>
                    ))}
                  </select>
                </FormControl>
                <ChevronDownIcon className='absolute end-3 top-2.5 h-4 w-4 opacity-50' />
              </div>
              <FormDescription className='font-manrope'>
                Set the font you want to use in the dashboard.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='theme'
          render={({ field }) => (
            <FormItem>
              <FormLabel>Theme</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger className='w-50'>
                    <SelectValue placeholder='Select a theme' />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value='light'>Light</SelectItem>
                  <SelectItem value='dark'>Dark</SelectItem>
                  <SelectItem value='violet'>Violet</SelectItem>
                  <SelectItem value='red'>Red</SelectItem>
                  <SelectItem value='rose'>Rose</SelectItem>
                  <SelectItem value='sky'>Sky</SelectItem>
                  <SelectItem value='teal'>Teal</SelectItem>
                  <SelectItem value='yellow'>Yellow</SelectItem>
                  <SelectItem value='amber'>Amber</SelectItem>
                  <SelectItem value='blue'>Blue</SelectItem>
                  <SelectItem value='cyan'>Cyan</SelectItem>
                  <SelectItem value='emerald'>Emerald</SelectItem>
                  <SelectItem value='fuchsia'>Fuchsia</SelectItem>
                  <SelectItem value='green'>Green</SelectItem>
                  <SelectItem value='indigo'>Indigo</SelectItem>
                  <SelectItem value='lime'>Lime</SelectItem>
                  <SelectItem value='orange'>Orange</SelectItem>
                  <SelectItem value='pink'>Pink</SelectItem>
                  <SelectItem value='purple'>Purple</SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                Select the color theme for the dashboard.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type='submit'>Update preferences</Button>
      </form>
    </Form>
  )
}
