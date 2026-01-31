# React Context 修复指南

## 问题描述

原项目中使用自定义 Hook `useAuth()` 来管理认证状态，但这导致了一个严重的问题：

**每个组件调用 `useAuth()` 时都会创建一个新的状态实例**，导致状态无法在不同组件间共享。

### 问题表现
- 登录页面调用 `login()` 后，`isAuthenticated` 更新了
- 但 App 组件中的 `isAuthenticated` 还是 `false`
- 导致登录成功后无法跳转到主页

### 原因分析

```typescript
// ❌ 错误的做法
export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(...);
  // 每次调用都创建新状态
  return { isAuthenticated, login, logout };
}

// LoginPage 和 App 各自调用 useAuth()
// 它们拥有不同的状态实例，无法共享
```

## 解决方案：使用 React Context

创建全局认证状态，所有组件共享同一个状态实例。

### 1. 创建 AuthContext

```typescript
// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, useCallback } from 'react';

interface AuthContextType {
  isAuthenticated: boolean;
  loading: boolean;
  user: any;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(true);
  // ... 状态管理

  return (
    <AuthContext.Provider value={{ isAuthenticated, loading, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
```

### 2. 在应用中使用 AuthProvider

```typescript
// src/App.tsx
import { AuthProvider } from './contexts/AuthContext';

const App = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </BrowserRouter>
  );
};
```

### 3. 在组件中使用

```typescript
// 任何组件中都可以使用
import { useAuth } from '../../contexts/AuthContext';

const LoginPage = () => {
  const { login, isAuthenticated } = useAuth();
  // 现在所有组件共享同一个状态
};
```

## 优势

✅ **全局状态共享**
- 所有组件访问同一个认证状态
- 状态变化立即反映到所有使用该状态的组件

✅ **避免 prop drilling**
- 不需要层层传递认证状态
- 任何组件都可以通过 `useAuth()` 访问

✅ **类型安全**
- TypeScript 类型检查
- 忘记包裹 Provider 会抛出错误

✅ **性能优化**
- 使用 `useCallback` 缓存函数
- 只在状态变化时重新渲染

## 迁移步骤

1. ✅ 创建 `AuthContext.tsx`
2. ✅ 在 `App.tsx` 中包裹 `AuthProvider`
3. ✅ 删除旧的 `hooks/useAuth.ts`
4. ✅ 更新所有导入语句

## 相关文件

- `src/contexts/AuthContext.tsx` - 认证上下文
- `src/App.tsx` - 应用入口（使用 AuthProvider）
- `src/pages/Login/index.tsx` - 登录页面（使用 useAuth）
- `src/components/Layout/MainLayout.tsx` - 主布局（使用 useAuth）

## 最佳实践

### 何时使用 Context

✅ **适合使用 Context 的场景：**
- 全局认证状态
- 主题设置
- 用户偏好
- 语言/国际化设置

❌ **不适合使用 Context 的场景：**
- 频繁变化的状态（如鼠标位置）
- 组件特定的局部状态
- 大型列表数据（使用专门的 state 管理）

### Context 性能优化

1. **拆分 Context**
   ```typescript
   // 不要把所有状态放在一个 Context
   // 按功能拆分成多个 Context
   <AuthProvider>
     <ThemeProvider>
       <App />
     </ThemeProvider>
   </AuthProvider>
   ```

2. **使用 useMemo/useCallback**
   ```typescript
   const value = useMemo(() => ({
     user,
     login,
     logout
   }), [user]);
   ```

3. **按需渲染**
   ```typescript
   // 将状态和状态更新函数分开
   <AuthStateContext.Provider value={state}>
     <AuthDispatchContext.Provider value={dispatch}>
       {children}
     </AuthDispatchContext.Provider>
   </AuthStateContext.Provider>
   ```

## 参考资料

- [React Context 官方文档](https://react.dev/reference/react/useContext)
- [How to use React Context effectively](https://kentcdodds.com/blog/how-to-use-react-context-effectively)
- [React Hooks: useState vs useContext](https://blog.logrocket.com/react-hooks-state-vs-context/)
