# 快速启动指南

## 首次启动

### 1. 确保环境准备就绪

```bash
# 检查 Go 版本
go version

# 检查 Node.js 版本
node -v

# 检查 pnpm 版本
pnpm -v

# 检查 MySQL 是否运行
mysql -u root -p -e "SELECT VERSION();"
```

### 2. 安装依赖

```bash
# 克隆或进入项目目录
cd file-manager-system

# 安装所有依赖（包括前端和后端）
make install
```

### 3. 初始化数据库

```bash
# 创建数据库和表
make migrate

# 或手动执行
mysql -u root -p < apps/api/sql/init.sql
```

### 4. 配置数据库连接

编辑 `apps/api/configs/config.yaml`：

```yaml
database:
  host: localhost
  port: 3306
  database: file_manager
  username: root
  password: "你的密码"  # 修改这里
```

### 5. 启动开发服务

```bash
# 同时启动前后端（推荐）
make dev

# 或分别启动
# 终端 1：启动后端
make dev:api

# 终端 2：启动前端
make dev:web
```

### 6. 访问应用

- 前端：http://localhost:3000
- 后端 API：http://localhost:8080
- 健康检查：http://localhost:8080/health

### 7. 注册账号

首次使用需要注册账号：

1. 访问 http://localhost:3000
2. 点击"没有账号？去注册"
3. 输入用户名和密码
4. 注册成功后登录

## 日常开发

### 启动服务

```bash
# 同时启动前后端
pnpm dev

# 或使用 Make
make dev
```

### 构建生产版本

```bash
# 构建前后端
make build

# 只构建前端
make build:web

# 只构建后端
make build:api
```

### 清理构建产物

```bash
make clean
```

## 开发工具

### VS Code

推荐安装的扩展：

- Go
- ESLint
- Prettier
- TypeScript
- Vite

### 浏览器插件

- React Developer Tools
- Redux DevTools（如果使用）

## 常见问题

### 端口被占用

**前端端口 3000 被占用**：

编辑 `apps/web/vite.config.ts`：

```typescript
server: {
  port: 3001,  // 修改为其他端口
  // ...
}
```

**后端端口 8080 被占用**：

编辑 `apps/api/configs/config.yaml`：

```yaml
server:
  port: 8081  # 修改为其他端口
```

同时更新 `apps/web/vite.config.ts` 中的代理配置：

```typescript
proxy: {
  '/api': {
    target: 'http://localhost:8081',  // 修改为新端口
    changeOrigin: true,
  },
}
```

### 数据库连接失败

1. 确认 MySQL 正在运行
2. 检查用户名和密码
3. 确认数据库已创建

### Token 过期

系统会自动刷新 Token，如果遇到 Token 过期：

1. 退出登录
2. 重新登录

### 上传失败

检查：

1. 文件大小是否超过限制（默认 5GB）
2. `uploads/` 目录是否有写入权限
3. 磁盘空间是否充足

## 调试技巧

### 后端调试

```bash
# 使用 Delve 调试器
cd apps/api
dlv debug cmd/server/main.go
```

### 前端调试

1. 打开浏览器开发者工具（F12）
2. 查看 Console 和 Network 标签
3. React DevTools 查看组件状态

### 查看日志

```bash
# 后端日志会直接输出到终端
make dev:api

# 或查看构建日志
tail -f apps/api/logs/server.log
```

## 下一步

- 阅读 [后端 API 文档](../apps/api/README.md)
- 阅读 [前端开发指南](../apps/web/README.md)
- 查看 [示例代码](../examples/)
