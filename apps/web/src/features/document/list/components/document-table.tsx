import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getDocumentList } from '@/api/document/list'
import { getApplicationOptions } from '@/api/document/application/options'
import type { Document } from '@/api/document/list'
import type { ApplicationOption } from '@/api/document/application/options'
import {
  ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'
import {
  File,
  FileText,
  Image as ImageIcon,
  Film,
  Music,
  Archive,
  FileCode,
  FileSpreadsheet,
  FileJson,
} from 'lucide-react'

// 文件类型图标映射
function getFileTypeIcon(fileName: string) {
  const ext = fileName.split('.').pop()?.toLowerCase() || ''

  const iconMap: Record<string, React.ComponentType<{ className?: string }>> = {
    // 图片
    jpg: ImageIcon,
    jpeg: ImageIcon,
    png: ImageIcon,
    gif: ImageIcon,
    bmp: ImageIcon,
    svg: ImageIcon,
    webp: ImageIcon,
    ico: ImageIcon,

    // 文档
    pdf: FileText,
    doc: FileText,
    docx: FileText,
    txt: FileText,
    rtf: FileText,
    odt: FileText,

    // 表格
    xls: FileSpreadsheet,
    xlsx: FileSpreadsheet,
    csv: FileSpreadsheet,
    ods: FileSpreadsheet,

    // 演示
    ppt: FileText,
    pptx: FileText,

    // 代码
    js: FileCode,
    jsx: FileCode,
    ts: FileCode,
    tsx: FileCode,
    html: FileCode,
    css: FileCode,
    json: FileJson,
    xml: FileCode,
    py: FileCode,
    java: FileCode,
    go: FileCode,
    rs: FileCode,
    c: FileCode,
    cpp: FileCode,
    h: FileCode,
    md: FileCode,

    // 视频
    mp4: Film,
    avi: Film,
    mkv: Film,
    mov: Film,
    wmv: Film,
    flv: Film,
    webm: Film,

    // 音频
    mp3: Music,
    wav: Music,
    flac: Music,
    aac: Music,
    ogg: Music,
    m4a: Music,

    // 压缩包
    zip: Archive,
    rar: Archive,
    '7z': Archive,
    tar: Archive,
    gz: Archive,
  }

  const Icon = iconMap[ext] || File
  return Icon
}

// 文件类型颜色映射
function getFileTypeColor(ext: string) {
  const colorMap: Record<string, string> = {
    // 图片 - 蓝色
    jpg: 'text-blue-500',
    jpeg: 'text-blue-500',
    png: 'text-blue-500',
    gif: 'text-blue-500',
    svg: 'text-blue-500',
    webp: 'text-blue-500',

    // 文档 - 红色
    pdf: 'text-red-500',
    doc: 'text-red-500',
    docx: 'text-red-500',
    txt: 'text-red-500',
    rtf: 'text-red-500',

    // 表格 - 绿色
    xls: 'text-green-500',
    xlsx: 'text-green-500',
    csv: 'text-green-500',

    // 演示 - 橙色
    ppt: 'text-orange-500',
    pptx: 'text-orange-500',

    // 代码 - 紫色
    js: 'text-purple-500',
    jsx: 'text-purple-500',
    ts: 'text-purple-500',
    tsx: 'text-purple-500',
    html: 'text-purple-500',
    css: 'text-purple-500',
    json: 'text-purple-500',
    xml: 'text-purple-500',
    py: 'text-purple-500',
    md: 'text-purple-500',

    // 视频 - 粉色
    mp4: 'text-pink-500',
    avi: 'text-pink-500',
    mkv: 'text-pink-500',
    mov: 'text-pink-500',

    // 音频 - 青色
    mp3: 'text-cyan-500',
    wav: 'text-cyan-500',
    flac: 'text-cyan-500',

    // 压缩包 - 黄色
    zip: 'text-yellow-500',
    rar: 'text-yellow-500',
    '7z': 'text-yellow-500',
  }

  return colorMap[ext.toLowerCase()] || 'text-gray-500'
}

const columns: ColumnDef<Document>[] = [
  {
    accessorKey: 'file_name',
    header: '文件名',
    cell: ({ row }) => {
      const fileName = row.getValue('file_name') as string
      const Icon = getFileTypeIcon(fileName)
      const ext = fileName.split('.').pop()?.toLowerCase() || ''
      const color = getFileTypeColor(ext)

      return (
        <div className="flex items-center gap-2">
          <Icon className={cn('h-4 w-4', color)} />
          <span className="font-medium">{fileName}</span>
        </div>
      )
    },
  },
  {
    accessorKey: 'file_size',
    header: '文件大小',
    cell: ({ row }) => {
      const size = row.getValue('file_size') as number
      const formatted = formatFileSize(size)
      return <div>{formatted}</div>
    },
  },
  {
    accessorKey: 'file_type',
    header: '文件类型',
    cell: ({ row }) => {
      const fileName = row.getValue('file_name') as string
      const ext = fileName.split('.').pop()?.toUpperCase() || '未知'
      const color = getFileTypeColor(ext.toLowerCase())
      const Icon = getFileTypeIcon(fileName)

      return (
        <div className="flex items-center gap-1.5">
          <Icon className={cn('h-3.5 w-3.5', color)} />
          <span className="capitalize text-xs font-medium">{ext}</span>
        </div>
      )
    },
  },
  {
    accessorKey: 'uploaded_by',
    header: '上传者类型',
    cell: ({ row }) => {
      const type = row.getValue('uploaded_by') as string
      return (
        <div className="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium">
          {type === 'user' ? '用户' : '应用'}
        </div>
      )
    },
  },
  {
    accessorKey: 'created_at',
    header: '上传时间',
    cell: ({ row }) => {
      const date = new Date(row.getValue('created_at') as string)
      return <div>{formatDate(date)}</div>
    },
  },
]

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

function formatDate(date: Date): string {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

export function DocumentTable() {
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
  const [keyword, setKeyword] = useState('')
  const [appIdentifier, setAppIdentifier] = useState<string>('')

  // 获取应用选项列表（不含敏感信息）
  const { data: applicationOptions = [] } = useQuery({
    queryKey: ['application-options'],
    queryFn: () => getApplicationOptions(),
  })

  // 获取文档列表
  const { data, isLoading } = useQuery({
    queryKey: ['documents', pagination.pageIndex + 1, pagination.pageSize, keyword, appIdentifier],
    queryFn: () =>
      getDocumentList({
        page: pagination.pageIndex + 1,
        page_size: pagination.pageSize,
        keyword,
        app_identifier: appIdentifier,
      }),
  })

  const documents = data?.documents || []
  const total = data?.total || 0

  const table = useReactTable({
    data: documents,
    columns,
    pageCount: Math.ceil(total / pagination.pageSize),
    state: { pagination },
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
  })

  return (
    <div className="space-y-4">
      {/* 筛选和搜索 */}
      <div className="flex gap-4">
        <div className="flex-1">
          <Input
            placeholder="搜索文件名..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            className="max-w-sm"
          />
        </div>
        <Select value={appIdentifier || 'all'} onValueChange={(value) => setAppIdentifier(value === 'all' ? '' : value)}>
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="选择应用标识" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部应用</SelectItem>
            {applicationOptions
              .filter((app: ApplicationOption) => app.app_identifier || app.app_account)
              .map((app: ApplicationOption) => {
                // 如果 app_identifier 为空，使用 app_account 作为 value
                const value = app.app_identifier || app.app_account
                // const displayIdentifier = app.app_identifier || `${app.app_account} (无标识)`

                return (
                  <SelectItem key={app.id} value={value}>
                    {app.app_name}
                  </SelectItem>
                )
              })}
          </SelectContent>
        </Select>
      </div>

      {/* 表格 */}
      <div className="rounded-md border">
        <div className="relative w-full overflow-auto">
          <table className="w-full caption-bottom text-sm">
            <thead className="[&_tr]:border-b">
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <th
                      key={header.id}
                      className="h-12 px-4 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0"
                    >
                      {header.isPlaceholder
                        ? null
                        : flexRender(
                            header.column.columnDef.header,
                            header.getContext()
                          )}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody className="[&_tr:last-child]:border-0">
              {isLoading ? (
                <tr>
                  <td
                    colSpan={columns.length}
                    className="p-4 text-center text-muted-foreground"
                  >
                    加载中...
                  </td>
                </tr>
              ) : documents.length === 0 ? (
                <tr>
                  <td
                    colSpan={columns.length}
                    className="p-4 text-center text-muted-foreground"
                  >
                    暂无数据
                  </td>
                </tr>
              ) : (
                table.getRowModel().rows?.map((row) => (
                  <tr
                    key={row.id}
                    className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted"
                  >
                    {row.getVisibleCells().map((cell) => (
                      <td
                        key={cell.id}
                        className="p-4 align-middle [&:has([role=checkbox])]:pr-0"
                      >
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* 分页 */}
      <div className="flex items-center justify-between px-2">
        <div className="flex-1 text-sm text-muted-foreground">
          共 {total} 条记录
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            上一页
          </Button>
          <div className="flex items-center gap-1">
            <span className="text-sm">
              第 {pagination.pageIndex + 1} 页，共 {table.getPageCount()} 页
            </span>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            下一页
          </Button>
        </div>
      </div>
    </div>
  )
}
