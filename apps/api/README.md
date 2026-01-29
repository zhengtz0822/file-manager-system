# 文档管理系统

基于 Go + Gin 开发的文档管理系统，支持大文件分片上传、下载和预览功能。

## 功能特性

- ✅ 用户认证（JWT）
- ✅ Token 黑名单机制（支持撤销）
- ✅ 注销登录
- ✅ 大文件分片上传
- ✅ 文档下载
- ✅ 文档预览（图片、PDF）
- ✅ 文档管理（列表、删除）
- ✅ UUID 作为文档标识
- ✅ MySQL 数据持久化
- ✅ 本地文件存储

## 技术栈

- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT (golang-jwt/jwt)
- **UUID**: Google UUID
- **密码加密**: bcrypt

## 项目结构

```
file-manager-service/
├── cmd/
│   └── server/
│       └── main.go              # 主程序入口
├── internal/
│   ├── handler/                 # HTTP 处理器
│   │   ├── auth.go
│   │   └── document.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth.go
│   │   └── document.go
│   ├── repository/              # 数据访问层
│   │   ├── user.go
│   │   ├── document.go
│   │   └── chunk.go
│   ├── middleware/              # 中间件
│   │   └── auth.go
│   ├── model/                   # 数据模型
│   │   ├── user.go
│   │   ├── document.go
│   │   ├── chunk.go
│   │   └── db.go
│   ├── pkg/                     # 工具包
│   │   ├── jwt/
│   │   ├── uuid/
│   │   └── storage/
│   ├── config/                  # 配置
│   │   └── config.go
│   └── router/                  # 路由
│       └── router.go
├── configs/
│   └── config.yaml              # 配置文件
├── sql/
│   └── init.sql                 # 数据库初始化脚本
├── uploads/                     # 文件存储目录
│   ├── chunks/                  # 分片临时目录
│   └── documents/               # 完整文件目录
├── go.mod
├── go.sum
└── README.md
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- MySQL 8.0+

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

创建数据库并导入初始化脚本：

```bash
mysql -u root -p < sql/init.sql
```

### 4. 修改配置

编辑 [configs/config.yaml](configs/config.yaml)，修改数据库连接信息：

```yaml
database:
  host: localhost
  port: 3306
  database: file_manager
  username: root
  password: "your_password"
```

### 5. 启动服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

### 6. 创建管理员用户

使用 API 注册用户或直接插入数据库：

```bash
# 使用 bcrypt 生成密码哈希
# 然后插入数据库
```

## API 文档

### 认证相关

#### 注册用户
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

#### 用户登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

Response:
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin"
    }
  }
}
```

#### 注销登录
```
POST /api/v1/auth/logout
Authorization: Bearer {token}

Response:
{
  "code": 200,
  "message": "注销成功"
}
```

### 文档上传（分片上传）

#### 1. 初始化上传
```
POST /api/v1/documents/chunks/init
Authorization: Bearer {token}
Content-Type: application/json

{
  "file_name": "example.pdf",
  "file_size": 10485760,
  "chunk_size": 5242880
}

Response:
{
  "code": 200,
  "message": "初始化成功",
  "data": {
    "upload_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_chunks": 2,
    "chunk_size": 5242880,
    "file_size": 10485760
  }
}
```

#### 2. 上传分片
```
POST /api/v1/documents/chunks/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data

upload_id: 550e8400-e29b-41d4-a716-446655440000
chunk_number: 1
file: [binary data]
```

#### 3. 完成上传
```
POST /api/v1/documents/chunks/complete
Authorization: Bearer {token}
Content-Type: application/json

{
  "upload_id": "550e8400-e29b-41d4-a716-446655440000"
}

Response:
{
  "code": 200,
  "message": "上传完成",
  "data": {
    "document_id": "660e8400-e29b-41d4-a716-446655440000",
    "file_name": "example.pdf",
    "file_size": 10485760
  }
}
```

### 文档管理

#### 获取文档列表
```
GET /api/v1/documents?page=1&page_size=10&keyword=test
Authorization: Bearer {token}
```

#### 获取文档详情
```
GET /api/v1/documents/{document_id}
Authorization: Bearer {token}
```

#### 删除文档
```
DELETE /api/v1/documents/{document_id}
Authorization: Bearer {token}
```

### 文档访问

#### 下载文档
```
GET /api/v1/documents/{document_id}/download
Authorization: Bearer {token}
```

#### 预览文档
```
GET /api/v1/documents/{document_id}/preview
Authorization: Bearer {token}
```

## 支持的文件类型

- 图片: jpg, jpeg, png, gif
- 文档: pdf, doc, docx, xls, xlsx, ppt, pptx
- 文本: txt, md

## 配置说明

[config.yaml](configs/config.yaml) 主要配置项：

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # 运行模式: debug/release

storage:
  max_file_size: 5368709120  # 最大文件大小 (5GB)
  chunk_size: 5242880        # 分片大小 (5MB)

jwt:
  secret: "your-secret-key"  # JWT 密钥
  expire_hours: 24           # Token 过期时间
```

## 开发建议

### 为其他系统集成

1. **认证方式**：使用 JWT Token，在请求头中携带 `Authorization: Bearer {token}`
2. **文档标识**：所有文档使用 UUID 作为唯一标识
3. **分片上传**：大文件建议使用分片上传，提高上传成功率
4. **错误处理**：统一返回格式，包含 code、message、data 字段

### 扩展建议

- 添加文件去重功能（基于 MD5）
- 添加文档分类/标签功能
- 添加访问日志记录
- 支持云存储（OSS、S3）
- 添加文档版本管理
- 添加病毒扫描功能

## License

MIT
