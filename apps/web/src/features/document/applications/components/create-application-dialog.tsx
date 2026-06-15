import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Plus, Loader2 } from 'lucide-react'

export function CreateApplicationDialog() {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [appName, setAppName] = useState('')
  const [appIdentifier, setAppIdentifier] = useState('')
  const [identifierError, setIdentifierError] = useState('')

  // 验证应用标识格式
  const validateIdentifier = (value: string): string => {
    if (!value) {
      return '应用标识不能为空'
    }
    if (value.length > 16) {
      return '应用标识最长16位'
    }
    if (!/^[a-zA-Z0-9]+$/.test(value)) {
      return '应用标识只能包含英文和数字'
    }
    return ''
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // 验证应用名称
    if (!appName.trim()) {
      alert('请输入应用名称')
      return
    }

    // 验证应用标识
    const error = validateIdentifier(appIdentifier.trim())
    if (error) {
      setIdentifierError(error)
      return
    }

    setLoading(true)
    try {
      const { createApplication } = await import('@/api/document/application')
      const result = await createApplication({
        app_name: appName.trim(),
        app_identifier: appIdentifier.trim(),
      })

      // 使用更友好的方式显示成功消息
      alert(`应用创建成功！\n\n应用账号: ${result.app_account}\n应用密钥: ${result.app_secret}\n\n请妥善保管密钥信息`)

      setOpen(false)
      setAppName('')
      setAppIdentifier('')
      setIdentifierError('')

      // 刷新列表
      window.location.reload()
    } catch (error: any) {
      const errorMsg = error.response?.data?.message || '创建失败'
      if (errorMsg.includes('应用标识') || errorMsg.includes('已存在')) {
        setIdentifierError(errorMsg)
      } else {
        alert(errorMsg)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleIdentifierChange = (value: string) => {
    setAppIdentifier(value)
    // 实时验证
    if (value) {
      const error = validateIdentifier(value)
      setIdentifierError(error)
    } else {
      setIdentifierError('')
    }
  }

  const handleDialogChange = (newOpen: boolean) => {
    setOpen(newOpen)
    if (!newOpen) {
      // 关闭时重置表单
      setAppName('')
      setAppIdentifier('')
      setIdentifierError('')
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleDialogChange}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          创建应用
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>创建应用</DialogTitle>
            <DialogDescription>
              创建新的应用，系统将自动生成应用账号和密钥用于接口认证。
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid gap-2">
              <Label htmlFor="app_name">应用名称</Label>
              <Input
                id="app_name"
                placeholder="请输入应用名称"
                value={appName}
                onChange={(e) => setAppName(e.target.value)}
                disabled={loading}
              />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="app_identifier">
                应用标识 <span className="text-destructive">*</span>
              </Label>
              <Input
                id="app_identifier"
                placeholder="请输入应用标识（英文数字，最长16位）"
                value={appIdentifier}
                onChange={(e) => handleIdentifierChange(e.target.value)}
                disabled={loading}
                className={identifierError ? 'border-destructive' : ''}
              />
              {identifierError && (
                <p className="text-xs text-destructive">{identifierError}</p>
              )}
              <p className="text-xs text-muted-foreground">
                示例：myapp, v1_system, app2024
              </p>
            </div>
            <div className="rounded-md bg-muted p-3 text-sm text-muted-foreground">
              <p className="font-medium mb-2">说明：</p>
              <ul className="list-disc list-inside space-y-1">
                <li>应用标识用于文件存储路径，创建后不可修改</li>
                <li>应用账号和密钥将由系统自动生成</li>
                <li>创建成功后将显示生成的认证信息</li>
                <li>请妥善保管应用密钥，用于JWT认证</li>
              </ul>
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => handleDialogChange(false)}
              disabled={loading}
            >
              取消
            </Button>
            <Button type="submit" disabled={loading}>
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              创建
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
