# 服务间调用指南

本文档说明如何作为子系统调用文档管理服务的 API。

## 认证方式

当前使用 **JWT Bearer Token** 认证方式，适合小规模服务间调用场景。

### 获取 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

响应示例：
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "your_username"
    }
  }
}
```

### 使用 Token

在所有需要认证的请求中添加 Header：
```
Authorization: Bearer {token}
```

示例：
```bash
curl -X GET http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token 管理

**Token 有效期**：默认 24 小时（可在配置文件中修改）

**Token 撤销**：
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer {token}"
```

## 推荐的集成方式

### 方案一：集中管理 Token（推荐）

在调用方系统中：
1. 创建一个配置服务或配置项存储 Token
2. 实现 Token 自动刷新逻辑
3. 所有子系统共享这个 Token

```go
// 示例：Go 客户端
type DocumentServiceClient struct {
    baseURL    string
    token      string
    httpClient *http.Client
}

func NewDocumentServiceClient(baseURL, username, password string) (*DocumentServiceClient, error) {
    client := &DocumentServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{},
    }

    // 登录获取 Token
    if err := client.login(username, password); err != nil {
        return nil, err
    }

    // 启动 Token 刷新协程
    go client.refreshTokenPeriodically()

    return client, nil
}

func (c *DocumentServiceClient) login(username, password string) error {
    // 实现登录逻辑
    return nil
}

func (c *DocumentServiceClient) refreshTokenPeriodically() {
    ticker := time.NewTicker(23 * time.Hour) // 在过期前刷新
    defer ticker.Stop()

    for range ticker.C {
        c.login(c.username, c.password)
    }
}
```

### 方案二：每个服务独立管理 Token

每个子系统维护自己的 Token，适用于需要独立授权的场景。

## 最佳实践

### 1. 安全存储 Token

```go
// 使用环境变量或配置管理系统
token := os.Getenv("DOC_SERVICE_TOKEN")
```

### 2. 错误处理

```go
resp, err := http.NewRequest("GET", url, nil)
if err != nil {
    // 处理错误
}

resp.Header.Set("Authorization", "Bearer " + token)

httpResp, err := client.Do(resp)
if httpResp.StatusCode == 401 {
    // Token 失效，重新登录
    client.refreshToken()
    // 重试请求
}
```

### 3. 重试机制

```go
func (c *DocumentServiceClient) doWithRetry(req *http.Request) (*http.Response, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        resp, err := c.httpClient.Do(req)
        if err != nil {
            if i < maxRetries-1 {
                time.Sleep(time.Second * time.Duration(i+1))
                continue
            }
            return nil, err
        }

        if resp.StatusCode == 401 {
            // Token 失效，刷新并重试
            c.refreshToken()
            req.Header.Set("Authorization", "Bearer "+c.token)
            continue
        }

        return resp, nil
    }
    return nil, errors.New("max retries exceeded")
}
```

### 4. 监控和日志

```go
func (c *DocumentServiceClient) do(req *http.Request) (*http.Response, error) {
    start := time.Now()

    resp, err := c.httpClient.Do(req)
    duration := time.Since(start)

    // 记录调用日志
    log.Printf("API call: %s %s - Status: %d - Duration: %v",
        req.Method, req.URL.Path, resp.StatusCode, duration)

    return resp, err
}
```

## 文档上传流程

### 小文件（< 10MB）

可以直接使用单次上传（需要添加该接口）。

### 大文件（≥ 10MB）

使用分片上传：

```go
func UploadLargeFile(filePath string, chunkSize int) error {
    // 1. 获取文件信息
    file, _ := os.Open(filePath)
    fileInfo, _ := file.Stat()
    fileSize := fileInfo.Size()
    fileName := filepath.Base(filePath)

    // 2. 初始化上传
    initResp := client.InitUpload(fileName, fileSize, chunkSize)
    uploadID := initResp.UploadID

    // 3. 分片上传
    buffer := make([]byte, chunkSize)
    chunkNumber := 1

    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }

        // 上传分片
        client.UploadChunk(uploadID, chunkNumber, buffer[:n])
        chunkNumber++
    }

    // 4. 完成上传
    completeResp := client.CompleteUpload(uploadID)

    return nil
}
```

## 常见问题

### Q: Token 过期了怎么办？

A: 实现自动刷新机制，在收到 401 响应时重新登录获取新 Token。

### Q: 如何提高性能？

A:
1. 使用连接池（`http.Client` 默认支持）
2. 实现客户端缓存
3. 使用长连接

### Q: 如何处理并发上传？

A: 使用 goroutine 并发上传分片，但要注意控制并发数量：

```go
semaphore := make(chan struct{}, 5) // 最多5个并发

for i := 1; i <= totalChunks; i++ {
    semaphore <- struct{}{} // 获取信号量
    go func(chunkNum int) {
        defer func() { <-semaphore }() // 释放信号量
        client.UploadChunk(uploadID, chunkNum, data)
    }(i)
}
```

### Q: 是否需要 HTTPS？

A: 生产环境强烈建议使用 HTTPS，防止 Token 被窃取。

## 未来扩展

如果系统规模扩大，可以考虑：

1. **API Gateway**：统一管理所有服务的认证
2. **OAuth 2.0**：使用 Client Credentials 流程
3. **服务网格**：使用 Istio 等服务网格的 mTLS
4. **Token 服务**：独立的 Token 管理服务
