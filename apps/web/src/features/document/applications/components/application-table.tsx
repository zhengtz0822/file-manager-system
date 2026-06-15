import { useEffect, useState } from 'react'
import type { Application } from '@/api/document/application'

export function ApplicationTable() {
  const [applications, setApplications] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const pageSize = 10

  const fetchApplications = async () => {
    setLoading(true)
    try {
      // Dynamic import to avoid circular dependency
      const { getApplicationList } = await import('@/api/document/application')
      const response = await getApplicationList({ page, page_size: pageSize })
      setApplications(response.applications || [])
      setTotal(response.total || 0)
    } catch (error) {
      console.error('获取应用列表失败', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchApplications()
  }, [page])

  const handleDelete = async (id: number) => {
    if (!confirm('确认删除该应用？')) return
    try {
      const { deleteApplication } = await import('@/api/document/application')
      await deleteApplication(id)
      fetchApplications()
    } catch (error) {
      console.error('删除失败', error)
    }
  }

  const handleToggleStatus = async (id: number, currentStatus: number) => {
    try {
      const { updateApplicationStatus } = await import('@/api/document/application')
      const newStatus = currentStatus === 1 ? 0 : 1
      await updateApplicationStatus(id, newStatus)
      fetchApplications()
    } catch (error) {
      console.error('状态更新失败', error)
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    alert('已复制到剪贴板')
  }

  return (
    <div className="rounded-md border">
      <table className="w-full">
        <thead>
          <tr className="border-b bg-muted">
            <th className="p-4 text-left">ID</th>
            <th className="p-4 text-left">应用名称</th>
            <th className="p-4 text-left">应用账号</th>
            <th className="p-4 text-left">应用密钥</th>
            <th className="p-4 text-left">状态</th>
            <th className="p-4 text-left">创建时间</th>
            <th className="p-4 text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          {loading ? (
            <tr>
              <td colSpan={7} className="p-4 text-center">
                加载中...
              </td>
            </tr>
          ) : applications.length === 0 ? (
            <tr>
              <td colSpan={7} className="p-4 text-center">
                暂无数据
              </td>
            </tr>
          ) : (
            applications.map((app) => (
              <tr key={app.id} className="border-b">
                <td className="p-4">{app.id}</td>
                <td className="p-4">{app.app_name}</td>
                <td className="p-4 font-mono text-sm">
                  <div className="flex items-center gap-2">
                    <span className="truncate max-w-[150px]">
                      {app.app_account}
                    </span>
                    <button
                      type="button"
                      onClick={() => copyToClipboard(app.app_account)}
                      className="text-blue-500 hover:text-blue-700"
                    >
                      复制
                    </button>
                  </div>
                </td>
                <td className="p-4 font-mono text-sm">
                  <div className="flex items-center gap-2">
                    <span className="truncate max-w-[120px]">
                      {app.app_secret.slice(0, 8)}...
                    </span>
                    <button
                      type="button"
                      onClick={() => copyToClipboard(app.app_secret)}
                      className="text-blue-500 hover:text-blue-700"
                    >
                      复制
                    </button>
                  </div>
                </td>
                <td className="p-4">
                  <span
                    className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                      app.status === 1
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {app.status === 1 ? '启用' : '禁用'}
                  </span>
                </td>
                <td className="p-4">{app.created_at}</td>
                <td className="p-4 text-right">
                  <div className="flex justify-end gap-2">
                    <button
                      type="button"
                      onClick={() =>
                        handleToggleStatus(app.id, app.status)
                      }
                      className="text-blue-500 hover:text-blue-700"
                    >
                      {app.status === 1 ? '禁用' : '启用'}
                    </button>
                    <button
                      type="button"
                      onClick={() => handleDelete(app.id)}
                      className="text-red-500 hover:text-red-700"
                    >
                      删除
                    </button>
                  </div>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>

      {/* 分页 */}
      {total > pageSize && (
        <div className="flex items-center justify-end gap-2 p-4">
          <button
            type="button"
            className="px-3 py-1 border rounded disabled:opacity-50"
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
          >
            上一页
          </button>
          <span className="text-sm text-muted-foreground">
            第 {page} 页，共 {Math.ceil(total / pageSize)} 页
          </span>
          <button
            type="button"
            className="px-3 py-1 border rounded disabled:opacity-50"
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(total / pageSize)}
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
