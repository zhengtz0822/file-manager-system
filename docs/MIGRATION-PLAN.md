# 前端项目迁移计划：从 Vite 到 Ant Design Pro (Umi)

## 📋 迁移概述

**目标**：将 `apps/web` 从手动搭建的 Vite 项目迁移到完整的 Ant Design Pro (Umi) 项目。

**保留**：在 `apps/web` 目录下，保持项目结构不变。

**迁移时间估计**：2-4 小时

---

## 🔄 迁移步骤

### 第一阶段：准备工作（30分钟）

#### 1.1 备份当前项目
```bash
# 在项目根目录执行
cd /Volumes/AppleExt/Codes/file-manager-system

# 创建备份
cp -r apps/web apps/web-backup-$(date +%Y%m%d)

# 或者使用 git（如果项目有版本控制）
git checkout -b backup-before-migration
git add .
git commit -m "备份：迁移到 Ant Design Pro 之前"
git push origin backup-before-migration
```

#### 1.2 清理当前 web 目录（可选）
```bash
# 如果确定完全重建，可以删除
# ⚠️ 警告：此操作不可逆，请确保已备份！
rm -rf apps/web/*

# 如果想保留某些文件（如文档），可以先移动
mv apps/web/README.md /tmp/
mv apps/web/docs /tmp/
```

#### 1.3 安装 Umi 脚手架工具
```bash
# 全局安装 umi（如果还没有）
npm install -g umi

# 或使用 pnpm
pnpm add -g umi
```

---

### 第二阶段：创建新项目（30分钟）

#### 2.1 使用 Umi 创建 Ant Design Pro 项目
```bash
# 进入 apps 目录
cd /Volumes/AppleExt/Codes/file-manager-system/apps

# 创建项目
umi create web

# 或使用完整命令（更多控制）
umi create web --template=ant-design-pro
```

#### 2.2 项目创建时的选择
在交互式命令行中选择：

```
? 选择模板
❯ Ant Design Pro

? TypeScript
❯ Yes

? 是否使用 @umijs/max
❯ Yes（推荐，新版 umi max）

? 是否使用 Ant Design Pro 组件
❯ Yes

? 是否需要布局
❯ Yes

? 是否需要权限
❯ Yes

? 是否需要国际化
❯ No（中文项目，暂时不需要）
❯ Yes（如果将来需要多语言）

? 是否需要代码规范（eslint/prettier）
❯ Yes

? 是否需要 mock 数据
❯ Yes

? 包管理工具
❯ pnpm
```

#### 2.3 安装依赖
```bash
cd web
pnpm install
```

#### 2.4 验证项目
```bash
# 启动开发服务器
pnpm dev

# 应该看到默认的 Ant Design Pro 页面
# 访问 http://localhost:8000
```

---

### 第三阶段：迁移代码（1-2小时）

#### 3.1 迁移 Context（全局状态管理）
```bash
# 源文件
apps/web-backup/src/contexts/AuthContext.tsx

# 目标位置
web/src/

# 操作：直接复制
cp apps/web-backup/src/contexts/AuthContext.tsx web/src/
```

**需要的修改**：将 React Router 的 `useNavigate` 改为 umi 的 `useNavigate`（如果 API 不同）

```typescript
// web/src/contexts/AuthContext.tsx
import { history } from '@umijs/max';

// 修改导航方式
navigate('/login');  // React Router
history.push('/login');  // Umi - 需要修改
```

#### 3.2 迁移 Services 层
```bash
# 复制整个 services 目录
cp -r apps/web-backup/src/services web/src/

# 检查并调整 API 基础路径
```

**修改 API 配置**（`web/src/services/api.ts`）：

```typescript
// 原代码
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
});

// 修改为 umi-request（Umi 推荐）
import { request } from '@umijs/max';

export const request = {
  config: {
    baseURL: '/api/v1',
    timeout: 10000,
  },
  requestInterceptors: [
    (url, options) => {
      const token = localStorage.getItem('token');
      if (token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${token}`,
        };
      }
      return { url, options };
    },
  ],
  responseInterceptors: [
    (response) => {
      const { code, data, message } = response as any;
      if (code === 200) {
        return data;
      }
      if (code === 401) {
        history.push('/login');
      }
      return response;
    },
  ],
};
```

#### 3.3 迁移页面组件
```bash
# 创建登录页面目录
mkdir -p web/src/pages/Login

# 复制登录页面
cp apps/web-backup/src/pages/Login/index.tsx web/src/pages/Login/
cp apps/web-backup/src/pages/Login/login.css web/src/pages/Login/

# 创建文档页面目录
mkdir -p web/src/pages/Document

# 复制文档页面
cp apps/web-backup/src/pages/Document/List.tsx web/src/pages/Document/
cp apps/web-backup/src/pages/Document/Upload.tsx web/src/pages/Document/
cp apps/web-backup/src/pages/Document/Preview.tsx web/src/pages/Document/
```

**需要修改的地方**：

1. **路由跳转**：
```typescript
// React Router
import { useNavigate } from 'react-router-dom';
const navigate = useNavigate();
navigate('/documents');

// Umi
import { history } from '@umijs/max';
history.push('/documents');
```

2. **获取路由参数**：
```typescript
// React Router
import { useParams } from 'react-router-dom';
const { id } = useParams();

// Umi
import { useParams } from '@umijs/max';
const { id } = useParams();
```

3. **当前路径**：
```typescript
// React Router
import { useLocation } from 'react-router-dom';
const location = useLocation();

// Umi
import { useLocation } from '@umijs/max';
const location = useLocation();
```

#### 3.4 迁移组件
```bash
# 创建 Layout 目录
mkdir -p web/src/components/Layout

# 复制 Ballpit 组件
cp apps/web-backup/src/components/Ballpit.jsx web/src/components/

# ⚠️ 注意：MainLayout 不需要复制
# Umi 会使用配置文件中的布局
```

#### 3.5 迁移类型定义
```bash
# 复制 types 目录
cp -r apps/web-backup/src/types web/src/
```

#### 3.6 迁移工具函数
```bash
# 复制 utils 目录
cp -r apps/web-backup/src/utils web/src/
```

#### 3.7 迁移 Hooks
```bash
# ⚠️ 注意：hooks/useAuth.ts 已经被 AuthContext 替代
# 不需要复制，已经在第 3.1 步处理
```

---

### 第四阶段：配置调整（30分钟）

#### 4.1 更新代理配置

**原配置** (`apps/web-backup/vite.config.ts`):
```typescript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
},
```

**新配置** (`web/.umirc.ts` 或 `web/config/config.ts`):

```typescript
import { defineConfig } from '@umijs/max';

export default defineConfig({
  // ...其他配置

  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      pathRewrite: {},
    },
  },
});
```

#### 4.2 配置路由

Umi 使用约定式路由，不需要手动配置。

**文件结构即路由**：
```
web/src/pages/
├── Login/
│   └── index.tsx          → /login
└── Document/
    ├── List.tsx           → /document/list
    ├── Upload.tsx         → /document/upload
    └── Preview.tsx        → /document/preview
```

**如果想自定义路由**，在 `.umirc.ts` 中配置：

```typescript
export default defineConfig({
  routes: [
    {
      path: '/login',
      component: 'Login',
    },
    {
      path: '/',
      component: '@/layouts/index',
      routes: [
        {
          path: '/documents',
          component: './Document/List',
        },
        {
          path: '/upload',
          component: './Document/Upload',
        },
        {
          path: '/documents/:id/preview',
          component: './Document/Preview',
        },
      ],
    },
  ],
});
```

#### 4.3 配置布局

**创建布局文件** (`web/src/layouts/index.tsx`):

```typescript
import { ProLayout } from '@ant-design/pro-components';
import { useNavigate } from '@umijs/max';

export default (props: any) => {
  return (
    <ProLayout
      title="文件管理系统"
      route={{
        path: '/',
        routes: [
          {
            path: '/documents',
            name: '文档列表',
            icon: <FileTextOutlined />,
          },
          {
            path: '/upload',
            name: '上传文档',
            icon: <UploadOutlined />,
          },
        ],
      }}
      menuItemRender={(menuItemProps, defaultDom) => {
        return (
          <div onClick={() => props.history.push(menuItemProps.path || '/')}>
            {defaultDom}
          </div>
        );
      }}
      {...props}
    >
      {props.children}
    </ProLayout>
  );
};
```

#### 4.4 配置主题

在 `.umirc.ts` 中配置：

```typescript
export default defineConfig({
  theme: {
    token: {
      colorPrimary: '#1890ff',
    },
  },
  antd: {
    // 全局配置 Ant Design
    configs: [
      {
        message: {
          duration: 3,
        },
      },
    ],
  },
});
```

---

### 第五阶段：特殊处理（30分钟）

#### 5.1 登录页面背景动画

**保留 Ballpit 组件**：

```bash
# 已在 3.4 步复制
# web/src/components/Ballpit.jsx
```

**调整导入路径**：

```typescript
// web/src/pages/Login/index.tsx
import Ballpit from '@/components/Ballpit';
```

#### 5.2 Tailwind CSS 集成

如果继续使用 Tailwind CSS：

```bash
# 安装依赖
pnpm add -D tailwindcss @tailwindcss/vite

# 在 Umi 中使用 Tailwind 需要特殊配置
```

**创建 `web/tailwind.config.js`**:

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

**在 `web/src/global.less` 中导入**：

```less
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**或者**：考虑放弃 Tailwind，直接使用 Ant Design 的样式系统（推荐）。

#### 5.3 TypeScript 类型

**创建全局类型定义** (`web/src/types/global.d.ts`):

```typescript
declare module 'umi' {
  export interface IRoute {
    path: string;
    component?: string;
    routes?: IRoute[];
    [key: string]: any;
  }
}
```

---

### 第六阶段：测试和验证（30分钟）

#### 6.1 功能测试清单

- [ ] 登录功能
  - [ ] 页面正常显示
  - [ ] Ballpit 背景动画正常
  - [ ] 登录成功跳转

- [ ] 主布局
  - [ ] 侧边栏显示正常
  - [ ] 菜单导航正常
  - [ ] 用户信息显示
  - [ ] 退出登录

- [ ] 文档列表
  - [ ] 表格正常显示
  - [ ] 搜索功能
  - [ ] 分页功能

- [ ] 文档上传
  - [ ] 上传表单
  - [ ] 进度显示

- [ ] 文档预览
  - [ ] 预览功能

- [ ] API 请求
  - [ ] 拦截器正常
  - [ ] 错误处理
  - [ ] Token 管理

#### 6.2 样式检查

- [ ] 主题色正确
- [ ] 响应式布局
- [ ] 暗色主题切换

#### 6.3 性能检查

```bash
# 构建生产版本
pnpm build

# 检查构建输出
ls -lh dist/

# 本地预览生产构建
pnpm preview
```

---

### 第七阶段：清理和优化（可选，15分钟）

#### 7.1 删除不需要的文件

```bash
# 如果使用 Umi 约定式路由
rm -rf web/src/App.tsx
rm -rf web/src/components/Layout/MainLayout.tsx

# 如果不使用 Tailwind
rm web/tailwind.config.js
rm -f web/src/index.css  # 如果只包含 Tailwind
```

#### 7.2 更新文档

```bash
# 复制备份的文档
cp /tmp/README.md web/
cp -r /tmp/docs web/

# 更新 CLAUDE.md
# 编辑 web/README.md，说明使用 Umi 框架
```

---

## 📦 迁移前后对比

### 目录结构对比

**迁移前 (Vite)**:
```
apps/web/
├── src/
│   ├── App.tsx              # 路由配置
│   ├── main.tsx             # 入口
│   ├── contexts/            # 全局状态
│   ├── pages/               # 页面
│   ├── services/            # API
│   ├── components/          # 组件
│   ├── hooks/               # Hooks
│   ├── types/               # 类型
│   ├── utils/               # 工具
│   └── vite.config.ts       # Vite 配置
├── package.json
└── tsconfig.json
```

**迁移后 (Umi)**:
```
apps/web/
├── src/
│   ├── .umi/                # Umi 生成（不需要修改）
│   ├── .umi-production/     # Umi 生产构建
│   ├── app.tsx              # 应用运行时配置（可选）
│   ├── global.less          # 全局样式
│   ├── layouts/             # 布局
│   ├── pages/               # 页面（约定式路由）
│   ├── services/            # API（迁移）
│   ├── components/          # 组件（迁移）
│   ├── contexts/            # 全局状态（迁移）
│   ├── types/               # 类型（迁移）
│   └── utils/               # 工具（迁移）
├── config/
│   └── config.ts           # Umi 配置
├── .umirc.ts               # 或 config/config.ts
├── package.json
└── tsconfig.json
```

---

## ⚠️ 注意事项和常见问题

### 1. API 请求方式

**Umi 推荐**使用 `@umijs/max` 的 `request`：

```typescript
import { request } from '@umijs/max';

export async function login(data: LoginRequest) {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    data,
  });
}
```

### 2. 路由跳转

```typescript
// ❌ 不要用 React Router
import { useNavigate } from 'react-router-dom';

// ✅ 使用 Umi 的 history
import { history } from '@umijs/max';

history.push('/documents');
history.replace('/login');
```

### 3. 获取当前路由

```typescript
import { useLocation, useSearchParams } from '@umijs/max';

const location = useLocation();
const [searchParams] = useSearchParams();
```

### 4. 权限控制

Umi 内置权限系统，在 `access.ts` 中配置：

```typescript
export default (initialState: any) => {
  const { loginUser } = initialState;
  return {
    canAdmin: loginUser && loginUser.role === 'admin',
  };
};
```

### 5. 端口冲突

- Vite 默认端口：3000
- Umi 默认端口：8000

**如果想使用 3000 端口**：

```typescript
// .umirc.ts
export default defineConfig({
  devServer: {
    port: 3000,
  },
});
```

### 6. Mock 数据

Umi 的 Mock 数据文件：

```typescript
// mock/documents.ts
export default {
  '/api/documents': [
    { id: '1', name: '文档1' },
    { id: '2', name: '文档2' },
  ],
};
```

---

## 🎯 完整执行脚本

将以下命令保存为 `migrate-frontend.sh` 并执行：

```bash
#!/bin/bash

set -e  # 遇到错误立即退出

echo "=== 开始迁移前端项目 ==="

# 1. 备份
echo "步骤 1/7: 备份当前项目..."
cd /Volumes/AppleExt/Codes/file-manager-system
BACKUP_DATE=$(date +%Y%m%d)
cp -r apps/web "apps/web-backup-$BACKUP_DATE"
echo "✅ 备份完成: apps/web-backup-$BACKUP_DATE"

# 2. 清理旧项目
echo "步骤 2/7: 清理旧 web 目录..."
rm -rf apps/web/*
echo "✅ 清理完成"

# 3. 创建新项目
echo "步骤 3/7: 创建 Ant Design Pro 项目..."
cd apps
# 使用非交互模式创建
echo "请按照提示操作..."
umi create web
# 或使用: npx create-umi@latest web --template=ant-design-pro
echo "✅ 项目创建完成"

# 4. 安装依赖
echo "步骤 4/7: 安装依赖..."
cd web
pnpm install
echo "✅ 依赖安装完成"

# 5. 迁移代码
echo "步骤 5/7: 迁移代码..."
cp -r "../apps/web-backup-$BACKUP_DATE/src/contexts" src/
cp -r "../apps/web-backup-$BACKUP_DATE/src/services" src/
cp -r "../apps/web-backup-$BACKUP_DATE/src/components" src/ || true
cp -r "../apps/web-backup-$BACKUP_DATE/src/pages/Document" src/pages/
cp "../apps/web-backup-$BACKUP_DATE/src/pages/Login/index.tsx" src/pages/Login/
cp "../apps/web-backup-$BACKUP_DATE/src/pages/Login/login.css" src/pages/Login/
cp -r "../apps/web-backup-$BACKUP_DATE/src/types" src/
cp -r "../apps/web-backup-$BACKUP_DATE/src/utils" src/ || true
echo "✅ 代码迁移完成"

# 6. 配置调整
echo "步骤 6/7: 配置调整..."
# TODO: 自动化配置修改
echo "请手动检查并更新配置文件"
echo "  - .umirc.ts 或 config/config.ts"
echo "  - 代理配置"
echo "  - 路由配置"
echo "  - 布局配置"

# 7. 启动项目
echo "步骤 7/7: 启动开发服务器..."
pnpm dev

echo "=== 迁移完成！==="
echo "请访问 http://localhost:8000 查看应用"
echo "记得检查和调整配置文件"
```

---

## 📚 参考资料

- [Umi 官方文档](https://umijs.org/)
- [Ant Design Pro 文档](https://pro.ant.design/)
- [从 Vite 迁移到 Umi](https://umijs.org/docs/max/vite)
- [Umi 速查表](https://umijs.org/docs/max/cheatsheet)

---

## ✅ 迁移检查清单

完成后使用此清单验证：

- [ ] 项目可以在 8000 端口启动（或配置的端口）
- [ ] 登录页面正常显示，包括 Ballpit 动画
- [ ] 登录功能正常
- [ ] 登录后跳转到文档列表
- [ ] 侧边栏菜单正常显示
- [ ] 菜单导航正常工作
- [ ] 右上角显示用户信息
- [ ] 退出登录功能正常
- [ ] 文档列表页正常
- [ ] 文档上传页正常
- [ ] 文档预览页正常
- [ ] API 请求正常（检查 Network 标签）
- [ ] 主题设置面板正常（右下角）
- [ ] 暗色主题切换正常
- [ ] 生产构建成功

---

## 🆘 需要帮助？

如果迁移过程中遇到问题，请提供：

1. 错误信息或截图
2. 相关的配置文件内容
3. 执行的具体步骤

我会协助你解决！
