import { createFileRoute, redirect } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/auth-store'
import { AuthenticatedLayout } from '@/components/layout/authenticated-layout'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: () => {
    const token = useAuthStore.getState().auth.accessToken
    if (!token) {
      throw redirect({
        to: '/sign-in',
      })
    }
  },
  component: AuthenticatedLayout,
})
