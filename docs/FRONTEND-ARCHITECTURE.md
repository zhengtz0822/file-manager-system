# 前端架构规范

## 目录结构

```
apps/web/src/
├── components/         # 可复用组件
│   └── Layout/        # 布局组件
├── pages/             # 页面组件（按业务模块组织）
│   ├── Login/         # 登录页面
│   └── Document/      # 文档管理页面
│       ├── List.tsx   # 文档列表
│       ├── Upload.tsx # 文档上传
│       └── Preview.tsx # 文档预览
├── services/          # API 服务层（按业务模块组织）
│   ├── api.ts         # 基础 axios 配置和拦截器
│   ├── auth/          # 认证相关 API
│   │   └── index.ts   # 导出认证相关函数和类型
│   └── document/      # 文档相关 API
│       └── index.ts   # 导出文档相关函数和类型
├── hooks/             # 自定义 React Hooks
├── types/             # TypeScript 类型定义
├── utils/             # 工具函数
└── App.tsx            # 应用入口和路由配置
```

## 开发规范

### 1. 目录组织原则

#### Pages 目录
- 每个业务模块一个文件夹
- 文件夹名使用 PascalCase（如 `Document`、`User`）
- 页面组件使用具名文件（如 `List.tsx`、`Detail.tsx`）
- 示例：
  ```
  pages/
  ├── Login/
  │   ├── index.tsx       # 登录页面
  │   └── login.css       # 页面样式
  └── Document/
      ├── List.tsx        # 文档列表
      ├── Upload.tsx      # 文档上传
      └── Preview.tsx     # 文档预览
  ```

#### Services 目录
- **`api.ts`**：基础 axios 配置，包括：
  - 创建 axios 实例
  - 请求/响应拦截器
  - 错误处理
  - 统一导出 request 实例

- **业务模块文件夹**：每个业务模块一个文件夹，包含：
  - `index.ts`：导出该模块的所有 API 函数和类型定义
  - API 函数命名：`get{资源}`、`create{资源}`、`update{资源}`、`delete{资源}`
  - 示例：
    ```typescript
    // services/document/index.ts
    export async function getDocumentList() { ... }
    export async function uploadDocument() { ... }
    export async function deleteDocument() { ... }
    ```

#### 导入规则
- 从 services 导入时使用模块名（自动解析 index.ts）：
  ```typescript
  // ✅ 正确
  import { login } from '@/services/auth';
  import { getDocumentList } from '@/services/document';

  // ❌ 避免
  import { login } from '@/services/auth/index';
  ```

### 2. API 调用规范

#### 在组件中使用
```typescript
import { getDocumentList } from '@/services/document';

function DocumentList() {
  const [documents, setDocuments] = useState([]);

  useEffect(() => {
    getDocumentList().then(setDocuments);
  }, []);

  return <div>...</div>;
}
```

#### 在 Hooks 中使用
```typescript
// hooks/useDocuments.ts
import { getDocumentList, deleteDocument } from '@/services/document';

export function useDocuments() {
  const [documents, setDocuments] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetch = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getDocumentList();
      setDocuments(data);
    } finally {
      setLoading(false);
    }
  }, []);

  // ...
}
```

### 3. 类型定义规范

#### API 相关类型
- 定义在对应的 service 文件中（与 API 函数在同一文件）
- 使用 TypeScript 接口定义请求和响应类型
- 示例：
  ```typescript
  // services/auth/index.ts
  export interface LoginRequest {
    username: string;
    password: string;
  }

  export interface LoginResponse {
    token: string;
    user: User;
  }
  ```

#### 通用类型
- 复用的类型定义在 `types/index.ts`
- 使用 `export` 导出供其他模块使用

### 4. 组件开发规范

#### 页面组件
- 使用函数组件 + Hooks
- 组件文件名使用 PascalCase
- Props 类型定义在组件文件顶部

#### 布局组件
- 放置在 `components/Layout/` 目录
- 包含导航、侧边栏、页头等

#### 可复用组件
- 放置在 `components/` 目录
- 使用 PascalCase 命名文件夹
- 示例：`components/DataTable/`

### 5. 路由规范

- 使用 React Router v6
- 路由配置在 `App.tsx` 中
- 受保护路由使用 `useAuth` hook 检查认证状态
- 示例：
  ```typescript
  <Route
    path="/"
    element={
      isAuthenticated ? (
        <MainLayout>
          <Routes>
            <Route path="documents" element={<DocumentList />} />
          </Routes>
        </MainLayout>
      ) : (
        <Navigate to="/login" replace />
      )
    }
  />
  ```

### 6. 状态管理规范

#### 本地状态
- 使用 `useState` 管理组件本地状态
- 使用 `useContext` 管理跨组件状态（如认证信息）

#### 服务端状态
- 优先使用 Hooks 封装 API 调用
- 使用 `useEffect` 在组件挂载时获取数据
- 提供 loading 状态和错误处理

### 7. 样式规范

#### CSS 模块
- 页面特定样式放在同名 CSS 文件中
- 使用 BEM 命名规范（可选）
- 示例：
  ```
  pages/
  └── Login/
      ├── index.tsx
      └── login.css
  ```

#### Tailwind CSS
- 优先使用 Tailwind 工具类
- 复杂样式使用 CSS 模块

### 8. 错误处理规范

#### API 错误
- 在 `api.ts` 的响应拦截器中统一处理
- 使用 Ant Design 的 `message` 组件显示错误
- 401 错误自动跳转到登录页

#### 组件错误
- 使用 `try-catch` 捕获异步错误
- 向用户显示友好的错误信息

### 9. 代码风格

#### TypeScript
- 启用严格模式
- 避免使用 `any`，优先使用具体类型或 `unknown`
- 函数参数和返回值必须标注类型

#### 命名规范
- 组件：PascalCase（如 `DocumentList`）
- 函数/变量：camelCase（如 `getDocumentList`）
- 类型/接口：PascalCase（如 `LoginRequest`）
- 常量：UPPER_SNAKE_CASE（如 `API_BASE_URL`）

## 最佳实践

### ✅ 推荐做法

1. **按业务模块组织代码**
   - 相关的页面、API、类型放在一起
   - 易于查找和维护

2. **使用自定义 Hooks 封装逻辑**
   - 减少组件中的重复代码
   - 提高代码复用性

3. **统一 API 调用**
   - 所有 API 调用通过 services 层
   - 便于统一管理和修改

4. **使用 TypeScript 类型**
   - 充分利用类型检查
   - 提高代码可维护性

### ❌ 避免的做法

1. **在组件中直接调用 axios**
   ```typescript
   // ❌ 避免
   useEffect(() => {
     axios.get('/api/v1/documents').then(...);
   }, []);

   // ✅ 推荐
   import { getDocumentList } from '@/services/document';
   useEffect(() => {
     getDocumentList().then(...);
   }, []);
   ```

2. **将所有 API 放在一个文件中**
   - 按业务模块拆分更易维护

3. **使用 CSS-in-JS（除非必要）**
   - 优先使用 Tailwind 或 CSS 模块

## 新增功能开发流程

1. 在 `pages/` 创建页面组件
2. 在 `services/{模块}/` 创建 API 函数
3. 在 `App.tsx` 添加路由
4. 在组件中使用 API
5. 添加必要的类型定义
6. 编写样式（CSS 模块或 Tailwind）

## 参考资源

- [React 文档](https://react.dev/)
- [React Router 文档](https://reactrouter.com/)
- [Ant Design 文档](https://ant.design/)
- [Tailwind CSS 文档](https://tailwindcss.com/)
