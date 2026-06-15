import { Main } from '@/components/layout/main'
import { ApplicationTable } from './components/application-table'
import { CreateApplicationDialog } from './components/create-application-dialog'

export function Applications() {
  return (
    <Main className="flex flex-1 flex-col gap-4 sm:gap-6">
      <div className="flex flex-wrap items-end justify-between gap-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">应用管理</h2>
          <p className="text-muted-foreground">
            管理外部系统调用认证应用
          </p>
        </div>
        <CreateApplicationDialog />
      </div>
      <ApplicationTable />
    </Main>
  )
}
