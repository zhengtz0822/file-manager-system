# 文件管理系统 (Monorepo)

基于 Go + Gin + React + Ant Design Pro 的文件管理系统，支持大文件分片上传、下载和预览。

## 项目结构

```
file-manager-system/
├── apps/
│   ├── api/                  # Go 后端服务
│   │   ├── cmd/              # 主程序入口
│   │   ├── internal/         # 内部包
│   │   ├── configs/          # 配置文件
│   │   └── sql/              # 数据库脚本
│   └── web/                  # React 前端应用
│       ├── src/
│       │   ├── components/   # 组件
│       │   ├── pages/        # 页面
│       │   ├── services/     # API 服务
│       │   └── hooks/        # React Hooks
│       └── package.json
├── packages/                 # 共享包
│   └── types/                # TypeScript 类型定义
├── pnpm-workspace.yaml       # pnpm workspace 配置
└── package.json              # 根 package.json
```

## 功能特性

### 后端 (Go)
- ✅ JWT 认证 + Token 黑名单
- ✅ 大文件分片上传 (支持 5GB)
- ✅ 文档下载和预览
- ✅ 文档管理 (CRUD)
- ✅ MySQL 数据持久化

### 前端 (React + Ant Design Pro)
- ✅ 用户登录/注销
- ✅ 文档列表 (分页、搜索)
- ✅ 大文件上传 (分片、进度条)
- ✅ 文档预览 (图片、PDF)
- ✅ 文档下载
- ✅ 响应式布局

## 技术栈

### 后端
- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT (golang-jwt/jwt)
- **密码**: bcrypt

### 前端
- **框架**: React 18
- **UI**: Ant Design 5 + Ant Design Pro Components
- **构建**: Vite 5
- **路由**: React Router v6
- **状态**: React Hooks
- **HTTP**: Axios
- **语言**: TypeScript

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- pnpm 8+
- MySQL 8.0+

### 1. 安装依赖

```bash
# 安装前端依赖
pnpm install

# 安装后端依赖（进入 api 目录）
cd apps/api
go mod download
cd ../..
```

### 2. 配置数据库

创建数据库并导入初始化脚本：

```bash
mysql -u root -p < apps/api/sql/init.sql
```

### 3. 修改配置

编辑 `apps/api/configs/config.yaml`：

```yaml
database:
  host: localhost
  port: 3306
  database: file_manager
  username: root
  password: "your_password"
```

### 4. 启动开发服务

#### 方式一：同时启动前后端（推荐）

```bash
pnpm dev
```

这会同时启动：
- 后端服务：http://localhost:8080
- 前端服务：http://localhost:3000

#### 方式二：分别启动

```bash
# 启动后端
pnpm dev:api

# 启动前端
pnpm dev:web
```

### 5. 访问应用

打开浏览器访问：http://localhost:3000

默认账号（需要先注册）：
- 用户名：admin
- 密码：自定义注册

## 开发指南

### 前端开发

```bash
# 进入 web 目录
cd apps/web

# 开发模式
pnpm dev

# 构建
pnpm build

# 预览构建结果
pnpm preview
```

### 后端开发

```bash
# 进入 api 目录
cd apps/api

# 运行
go run cmd/server/main.go

# 构建
go build -o ../../bin/server cmd/server/main.go
```

### API 代理配置

前端开发环境通过 Vite proxy 代理后端 API，配置在 `apps/web/vite.config.ts`：

```typescript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

## 生产部署

### 构建前端

```bash
pnpm build:web
```

构建产物在 `apps/web/dist` 目录。

### 构建后端

```bash
pnpm build:api
```

二进制文件在 `bin/server`。

### 部署方式

#### 方式一：前后端分离部署

1. 前端部署到 Nginx/CDN
2. 后端部署为独立服务

#### 方式二：Go Embed 单文件部署

将前端构建产物嵌入 Go 二进制（待实现）。

## 常见问题

### Q: 前端无法连接后端？

A: 检查后端是否启动，端口是否正确（默认 8080）。

### Q: 上传大文件失败？

A: 检查配置中的 `max_file_size` 和 `chunk_size` 设置。

### Q: Token 过期怎么办？

A: 系统会自动刷新 Token，如果还是失败请重新登录。

## 相关文档

- [后端 API 文档](apps/api/README.md)
- [前端开发指南](apps/web/README.md)
- [服务间调用指南](apps/api/docs/SERVICE-TO-SERVICE.md)
- [集成指南](apps/api/docs/INTEGRATION-GUIDE.md)

## License

MIT
