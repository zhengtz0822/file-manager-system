import { createFileRoute } from '@tanstack/react-router'
import { DocumentList } from '@/features/document/list'

export const Route = createFileRoute('/_authenticated/document/list/')({
  component: DocumentList,
})
