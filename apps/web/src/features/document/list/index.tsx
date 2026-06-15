import { Main } from '@/components/layout/main'
import { DocumentTable } from './components/document-table'

export function DocumentList() {
  return (
    <Main className="flex flex-1 flex-col gap-4 sm:gap-6">
      <div className="flex flex-wrap items-end justify-between gap-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">文档管理</h2>
          <p className="text-muted-foreground">
            查看和管理已上传的文档
          </p>
        </div>
      </div>
      <DocumentTable />
    </Main>
  )
}
