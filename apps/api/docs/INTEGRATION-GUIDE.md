# 服务间调用完整流程指南

本文档详细说明其他服务如何调用文档管理系统的上传接口。

## 整体流程

```
┌─────────────────┐
│  其他服务 A     │
│  (调用方)       │
└────────┬────────┘
         │
         │ 1. 登录获取 Token
         ↓
┌─────────────────────────────┐
│ POST /api/v1/auth/login     │
│ {username, password}        │
└────────┬────────────────────┘
         │
         │ 返回 Token
         ↓
┌─────────────────────────────┐
│ 存储 Token                  │
│ - 内存变量                  │
│ - 环境变量                  │
│ - 配置中心                  │
└────────┬────────────────────┘
         │
         │ 2. 初始化上传
         ↓
┌─────────────────────────────┐
│ POST /api/v1/documents/     │
│      chunks/init            │
│ Authorization: Bearer {token}│
└────────┬────────────────────┘
         │
         │ 返回 upload_id
         ↓
┌─────────────────────────────┐
│ 3. 循环上传分片             │
│ POST /api/v1/documents/     │
│      chunks/upload          │
│ (每个分片携带 Token)        │
└────────┬────────────────────┘
         │
         │ 4. 完成上传
         ↓
┌─────────────────────────────┐
│ POST /api/v1/documents/     │
│      chunks/complete        │
│ Authorization: Bearer {token}│
└────────┬────────────────────┘
         │
         │ 返回 document_id
         ↓
    上传完成 ✅
```

## 方案一：集中式 Token 管理（推荐）

### 架构说明

在调用方系统中创建一个**统一的服务客户端**，所有子系统都通过这个客户端调用文档服务。

```
┌──────────────────┐
│  子系统 A        │
│  (业务系统)      │
└────────┬─────────┘
         │
         ↓ 调用
┌──────────────────┐
│  文档服务客户端  │  ← Token 管理在这里
│  (SDK/封装层)    │  ← 自动刷新
└────────┬─────────┘
         │ HTTP 请求
         ↓
┌──────────────────┐
│  文档管理系统    │
│  (API 服务)      │
└──────────────────┘
```

### Go 实现示例

#### 1. 创建服务客户端

```go
package docsclient

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "sync"
    "time"
)

// DocumentServiceClient 文档服务客户端
type DocumentServiceClient struct {
    baseURL     string
    username    string
    password    string
    token       string
    tokenExpiry time.Time
    httpClient  *http.Client
    mutex       sync.RWMutex
}

// Config 客户端配置
type Config struct {
    BaseURL  string        // 文档服务地址，如 "http://localhost:8080"
    Username string        // 用户名
    Password string        // 密码
    Timeout  time.Duration // 请求超时
}

// NewClient 创建客户端
func NewClient(cfg *Config) (*DocumentServiceClient, error) {
    client := &DocumentServiceClient{
        baseURL:  cfg.BaseURL,
        username: cfg.Username,
        password: cfg.Password,
        httpClient: &http.Client{
            Timeout: cfg.Timeout,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
    }

    // 初始化 Token
    if err := client.login(); err != nil {
        return nil, fmt.Errorf("初始化登录失败: %w", err)
    }

    // 启动 Token 自动刷新协程
    go client.refreshTokenPeriodically()

    return client, nil
}

// login 登录获取 Token
func (c *DocumentServiceClient) login() error {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    loginReq := map[string]string{
        "username": c.username,
        "password": c.password,
    }

    body, _ := json.Marshal(loginReq)
    req, err := http.NewRequest("POST", c.baseURL+"/api/v1/auth/login", bytes.NewBuffer(body))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("登录失败，状态码: %d", resp.StatusCode)
    }

    var result struct {
        Code    int    `json:"code"`
        Message string `json:"message"`
        Data    struct {
            Token string `json:"token"`
            User  struct {
                ID       int    `json:"id"`
                Username string `json:"username"`
            } `json:"user"`
        } `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }

    c.token = result.Data.Token
    c.tokenExpiry = time.Now().Add(23 * time.Hour) // 提前1小时刷新

    return nil
}

// refreshTokenPeriodically 定期刷新 Token
func (c *DocumentServiceClient) refreshTokenPeriodically() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        if time.Now().After(c.tokenExpiry) {
            if err := c.login(); err != nil {
                fmt.Printf("刷新 Token 失败: %v\n", err)
            }
        }
    }
}

// getAuthToken 获取有效的 Token
func (c *DocumentServiceClient) getAuthToken() string {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return c.token
}

// doRequest 执行 HTTP 请求（带重试）
func (c *DocumentServiceClient) doRequest(method, path string, body io.Reader) (*http.Response, error) {
    maxRetries := 3
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        req, err := http.NewRequest(method, c.baseURL+path, body)
        if err != nil {
            return nil, err
        }

        // 添加 Token
        req.Header.Set("Authorization", "Bearer "+c.getAuthToken())
        if method == "POST" || method == "PUT" {
            req.Header.Set("Content-Type", "application/json")
        }

        resp, err := c.httpClient.Do(req)
        if err != nil {
            lastErr = err
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }

        // 检查是否 Token 失效
        if resp.StatusCode == http.StatusUnauthorized {
            resp.Body.Close()

            // 刷新 Token 并重试
            if err := c.login(); err != nil {
                lastErr = err
                continue
            }

            // 重建请求并重试
            if body != nil && method != "POST" {
                // 对于有 body 的请求，需要重新创建
                // 这里简化处理，实际需要缓存 body
            }
            continue
        }

        return resp, nil
    }

    return nil, fmt.Errorf("请求失败，重试 %d 次后仍失败: %w", maxRetries, lastErr)
}

// InitUpload 初始化上传
func (c *DocumentServiceClient) InitUpload(fileName string, fileSize int64, chunkSize int) (*InitUploadResponse, error) {
    reqBody := map[string]interface{}{
        "file_name":  fileName,
        "file_size":  fileSize,
        "chunk_size": chunkSize,
    }

    jsonData, _ := json.Marshal(reqBody)
    resp, err := c.doRequest("POST", "/api/v1/documents/chunks/init", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Code    int                `json:"code"`
        Message string             `json:"message"`
        Data    InitUploadResponse `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result.Data, nil
}

// UploadChunk 上传分片
func (c *DocumentServiceClient) UploadChunk(uploadID string, chunkNumber int, data []byte, fileName string) error {
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    writer.WriteField("upload_id", uploadID)
    writer.WriteField("chunk_number", fmt.Sprintf("%d", chunkNumber))

    part, err := writer.CreateFormFile("file", fileName)
    if err != nil {
        return err
    }

    if _, err := part.Write(data); err != nil {
        return err
    }

    writer.Close()

    req, err := http.NewRequest("POST", c.baseURL+"/api/v1/documents/chunks/upload", body)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer "+c.getAuthToken())
    req.Header.Set("Content-Type", writer.FormDataContentType())

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("上传分片失败，状态码: %d", resp.StatusCode)
    }

    return nil
}

// CompleteUpload 完成上传
func (c *DocumentServiceClient) CompleteUpload(uploadID string) (*CompleteUploadResponse, error) {
    reqBody := map[string]string{
        "upload_id": uploadID,
    }

    jsonData, _ := json.Marshal(reqBody)
    resp, err := c.doRequest("POST", "/api/v1/documents/chunks/complete", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Code    int                    `json:"code"`
        Message string                 `json:"message"`
        Data    CompleteUploadResponse `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result.Data, nil
}

// UploadFile 上传文件（完整流程）
func (c *DocumentServiceClient) UploadFile(filePath string, chunkSize int) (*CompleteUploadResponse, error) {
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        return nil, err
    }

    fileName := fileInfo.Name()
    fileSize := fileInfo.Size()

    // 1. 初始化上传
    fmt.Printf("初始化上传: %s (%d bytes)\n", fileName, fileSize)
    initResp, err := c.InitUpload(fileName, fileSize, chunkSize)
    if err != nil {
        return nil, fmt.Errorf("初始化上传失败: %w", err)
    }

    // 2. 上传分片
    buffer := make([]byte, chunkSize)
    chunkNumber := 1
    totalChunks := initResp.TotalChunks

    for {
        n, err := file.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("读取文件失败: %w", err)
        }

        fmt.Printf("上传分片 %d/%d\n", chunkNumber, totalChunks)
        if err := c.UploadChunk(initResp.UploadID, chunkNumber, buffer[:n], fileName); err != nil {
            return nil, fmt.Errorf("上传分片 %d 失败: %w", chunkNumber, err)
        }

        chunkNumber++
    }

    // 3. 完成上传
    fmt.Println("完成上传...")
    completeResp, err := c.CompleteUpload(initResp.UploadID)
    if err != nil {
        return nil, fmt.Errorf("完成上传失败: %w", err)
    }

    fmt.Printf("上传成功！文档ID: %s\n", completeResp.DocumentID)
    return completeResp, nil
}

// 响应结构体
type InitUploadResponse struct {
    UploadID    string `json:"upload_id"`
    TotalChunks int    `json:"total_chunks"`
    ChunkSize   int    `json:"chunk_size"`
    FileSize    int64  `json:"file_size"`
}

type CompleteUploadResponse struct {
    DocumentID string `json:"document_id"`
    FileName   string `json:"file_name"`
    FileSize   int64  `json:"file_size"`
}
```

#### 2. 使用示例

```go
package main

import (
    "fmt"
    "time"
    "yourproject/docsclient"
)

func main() {
    // 创建客户端（会自动登录并获取 Token）
    client, err := docsclient.NewClient(&docsclient.Config{
        BaseURL:  "http://localhost:8080",
        Username: "your_username",
        Password: "your_password",
        Timeout:  30 * time.Second,
    })
    if err != nil {
        panic(err)
    }

    // 上传文件（自动处理分片）
    resp, err := client.UploadFile("/path/to/your/file.pdf", 5*1024*1024) // 5MB 分片
    if err != nil {
        panic(err)
    }

    fmt.Printf("上传成功！文档ID: %s\n", resp.DocumentID)
}
```

## 方案二：环境变量配置（简单场景）

### 1. 在配置文件中设置

```yaml
# config.yaml
document_service:
  base_url: "http://localhost:8080"
  token: "${DOC_SERVICE_TOKEN}"  # 从环境变量读取
```

### 2. 初始化脚本获取 Token

```bash
#!/bin/bash
# scripts/init-doc-service-token.sh

# 登录获取 Token
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }')

# 提取 Token
TOKEN=$(echo $RESPONSE | jq -r '.data.token')

# 导出为环境变量
export DOC_SERVICE_TOKEN=$TOKEN

echo "Token 已设置: ${TOKEN:0:20}..."
```

### 3. 使用 Token 调用接口

```go
package main

import (
    "net/http"
    "os"
)

func uploadDocument(filePath string) error {
    token := os.Getenv("DOC_SERVICE_TOKEN")

    // 直接使用 Token 调用
    req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/documents/chunks/init", body)
    req.Header.Set("Authorization", "Bearer "+token)

    // ...
    return nil
}
```

## 方案三：配置中心（企业级）

### 使用 Nacos/Consul 等配置中心

```
┌──────────────────┐
│  配置中心        │
│  (Nacos/Consul) │
│                  │
│  - doc_service   │
│    - token: xxx  │
│    - url: xxx    │
└────────┬─────────┘
         │
         │ 所有服务订阅配置
         ↓
┌──────────────────┐
│  子系统 A/B/C    │
└──────────────────┘
```

## Python 调用示例

```python
import requests
import os
from typing import Optional

class DocumentServiceClient:
    def __init__(self, base_url: str, username: str, password: str):
        self.base_url = base_url
        self.username = username
        self.password = password
        self.token: Optional[str] = None
        self.login()

    def login(self):
        """登录获取 Token"""
        resp = requests.post(
            f"{self.base_url}/api/v1/auth/login",
            json={
                "username": self.username,
                "password": self.password
            }
        )
        resp.raise_for_status()
        data = resp.json()
        self.token = data["data"]["token"]

    def _get_headers(self):
        """获取请求头"""
        if not self.token:
            self.login()
        return {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        }

    def upload_file(self, file_path: str, chunk_size: int = 5*1024*1024):
        """上传文件"""
        file_name = os.path.basename(file_path)
        file_size = os.path.getsize(file_path)

        # 1. 初始化上传
        resp = requests.post(
            f"{self.base_url}/api/v1/documents/chunks/init",
            headers=self._get_headers(),
            json={
                "file_name": file_name,
                "file_size": file_size,
                "chunk_size": chunk_size
            }
        )
        resp.raise_for_status()
        upload_data = resp.json()["data"]
        upload_id = upload_data["upload_id"]

        # 2. 上传分片
        with open(file_path, 'rb') as f:
            chunk_num = 1
            while True:
                chunk = f.read(chunk_size)
                if not chunk:
                    break

                print(f"上传分片 {chunk_num}/{upload_data['total_chunks']}")

                files = {'file': (file_name, chunk)}
                data = {
                    'upload_id': upload_id,
                    'chunk_number': str(chunk_num)
                }

                resp = requests.post(
                    f"{self.base_url}/api/v1/documents/chunks/upload",
                    headers={"Authorization": f"Bearer {self.token}"},
                    files=files,
                    data=data
                )
                resp.raise_for_status()
                chunk_num += 1

        # 3. 完成上传
        resp = requests.post(
            f"{self.base_url}/api/v1/documents/chunks/complete",
            headers=self._get_headers(),
            json={"upload_id": upload_id}
        )
        resp.raise_for_status()
        return resp.json()["data"]

# 使用示例
if __name__ == "__main__":
    client = DocumentServiceClient(
        base_url="http://localhost:8080",
        username="your_username",
        password="your_password"
    )

    result = client.upload_file("/path/to/file.pdf")
    print(f"上传成功！文档ID: {result['document_id']}")
```

## Java 调用示例

```java
import okhttp3.*;
import org.json.JSONObject;
import java.io.File;
import java.io.FileInputStream;
import java.util.concurrent.TimeUnit;

public class DocumentServiceClient {
    private final String baseUrl;
    private final String username;
    private final String password;
    private String token;
    private final OkHttpClient client;

    public DocumentServiceClient(String baseUrl, String username, String password) {
        this.baseUrl = baseUrl;
        this.username = username;
        this.password = password;
        this.client = new OkHttpClient.Builder()
            .connectTimeout(30, TimeUnit.SECONDS)
            .build();
        login();
    }

    private void login() {
        try {
            JSONObject json = new JSONObject();
            json.put("username", username);
            json.put("password", password);

            RequestBody body = RequestBody.create(
                json.toString(),
                MediaType.parse("application/json")
            );

            Request request = new Request.Builder()
                .url(baseUrl + "/api/v1/auth/login")
                .post(body)
                .build();

            try (Response response = client.newCall(request).execute()) {
                JSONObject resp = new JSONObject(response.body().string());
                this.token = resp.getJSONObject("data").getString("token");
            }
        } catch (Exception e) {
            throw new RuntimeException("登录失败", e);
        }
    }

    public String uploadFile(String filePath, int chunkSize) throws Exception {
        File file = new File(filePath);
        String fileName = file.getName();
        long fileSize = file.length();

        // 1. 初始化上传
        JSONObject initJson = new JSONObject();
        initJson.put("file_name", fileName);
        initJson.put("file_size", fileSize);
        initJson.put("chunk_size", chunkSize);

        RequestBody initBody = RequestBody.create(
            initJson.toString(),
            MediaType.parse("application/json")
        );

        Request initRequest = new Request.Builder()
            .url(baseUrl + "/api/v1/documents/chunks/init")
            .post(initBody)
            .addHeader("Authorization", "Bearer " + token)
            .build();

        String uploadId;
        try (Response response = client.newCall(initRequest).execute()) {
            JSONObject resp = new JSONObject(response.body().string());
            uploadId = resp.getJSONObject("data").getString("upload_id");
        }

        // 2. 上传分片
        FileInputStream fis = new FileInputStream(file);
        byte[] buffer = new byte[chunkSize];
        int chunkNumber = 1;
        int bytesRead;

        while ((bytesRead = fis.read(buffer)) != -1) {
            System.out.println("上传分片 " + chunkNumber);

            RequestBody requestBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("upload_id", uploadId)
                .addFormDataPart("chunk_number", String.valueOf(chunkNumber))
                .addFormDataPart("file", fileName,
                    RequestBody.create(
                        okhttp3.MediaType.parse("application/octet-stream"),
                        buffer,
                        0,
                        bytesRead
                    ))
                .build();

            Request request = new Request.Builder()
                .url(baseUrl + "/api/v1/documents/chunks/upload")
                .post(requestBody)
                .addHeader("Authorization", "Bearer " + token)
                .build();

            try (Response response = client.newCall(request).execute()) {
                if (!response.isSuccessful()) {
                    throw new RuntimeException("上传分片失败");
                }
            }

            chunkNumber++;
        }
        fis.close();

        // 3. 完成上传
        JSONObject completeJson = new JSONObject();
        completeJson.put("upload_id", uploadId);

        RequestBody completeBody = RequestBody.create(
            completeJson.toString(),
            MediaType.parse("application/json")
        );

        Request completeRequest = new Request.Builder()
            .url(baseUrl + "/api/v1/documents/chunks/complete")
            .post(completeBody)
            .addHeader("Authorization", "Bearer " + token)
            .build();

        try (Response response = client.newCall(completeRequest).execute()) {
            JSONObject resp = new JSONObject(response.body().string());
            return resp.getJSONObject("data").getString("document_id");
        }
    }
}
```

## 最佳实践建议

### 1. Token 管理

✅ **推荐**：
- 使用统一的服务客户端封装
- 实现自动刷新机制
- 使用连接池提高性能

❌ **不推荐**：
- 每次请求都重新登录
- 硬编码 Token
- 多处复制 Token 获取逻辑

### 2. 错误处理

```go
// 实现重试机制
func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
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

        // Token 失效，刷新后重试
        if resp.StatusCode == 401 {
            resp.Body.Close()
            c.login()
            req.Header.Set("Authorization", "Bearer "+c.token)
            continue
        }

        return resp, nil
    }
    return nil, errors.New("max retries exceeded")
}
```

### 3. 监控和日志

```go
// 添加调用日志
func (c *Client) logAPICall(method, path string, statusCode int, duration time.Duration) {
    log.Printf("[文档服务] %s %s - 状态码: %d - 耗时: %v",
        method, path, statusCode, duration)
}
```

## 总结

**推荐方案**：使用**方案一（服务客户端封装）**

**优势**：
- ✅ Token 自动管理
- ✅ 统一错误处理
- ✅ 自动重试机制
- ✅ 易于维护
- ✅ 性能优化（连接池）

**实施步骤**：
1. 创建服务客户端 SDK（上面的 Go 代码）
2. 在各子系统中引入这个 SDK
3. 配置用户名和密码
4. 直接调用 SDK 方法上传文件

这样其他服务只需要关心业务逻辑，不需要处理 Token 管理细节！
