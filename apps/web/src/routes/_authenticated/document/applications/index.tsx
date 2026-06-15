import { createFileRoute } from '@tanstack/react-router'
import { Applications } from '@/features/document/applications'

export const Route = createFileRoute('/_authenticated/document/applications/')({
  component: Applications,
})
