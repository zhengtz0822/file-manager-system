# 主题和布局配置指南

## 功能概述

本系统集成了 Ant Design Pro 的主题工具，提供以下功能：

- ✅ **主题切换**：亮色/暗色主题
- ✅ **主题色定制**：8 种预设主题色
- ✅ **布局模式**：侧边栏/顶栏/混合布局
- ✅ **内容宽度**：流式/定宽
- ✅ **固定顶栏**：可开启/关闭
- ✅ **固定侧边栏**：可开启/关闭
- ✅ **分组菜单**：可开启/关闭
- ✅ **设置持久化**：自动保存到浏览器

## 使用方法

### 1. 打开设置面板

登录系统后，点击右下角的 **⚙️ 设置按钮**，即可打开主题设置面板。

### 2. 可配置项

#### 主题 (Theme)
- **Light**：亮色主题（默认）
- **Dark**：暗色主题

#### 主题色 (Primary Color)
可选择 8 种预设主题色：
- 蓝色 (Daybreak) - 默认
- 红色 (Dust)
- 橙色 (Volcano)
- 黄色 (Sunset)
- 青色 (Cyan)
- 绿色 (Green)
- 深蓝 (Ming)
- 紫色 (Purple)

#### 布局模式 (Layout)
- **Side**：侧边栏布局（默认）
- **Top**：顶栏布局
- **Mix**：混合布局

#### 内容宽度 (Content Width)
- **Fluid**：流式布局（默认）
- **Fixed**：定宽布局

#### 固定顶栏 (Fixed Header)
- 开启后，顶栏固定在顶部，内容滚动时顶栏不动

#### 固定侧边栏 (Fix Siderbar)
- 开启后，侧边栏固定，内容滚动时侧边栏不动

#### 分组菜单 (Split Menus)
- 开启后，菜单自动分组显示

### 3. 快捷主题切换

在顶部用户菜单中点击 **"设置"**，可以快速切换亮色/暗色主题。

## 技术实现

### 核心组件

使用 `@ant-design/pro-components` 的 `ProLayout` 组件：

```tsx
import { ProLayout } from '@ant-design/pro-components';

<ProLayout
  title="文件管理系统"
  theme={settings.theme}
  settings={settings}
  onSettingChange={handleSettingChange}
>
  <Outlet />
</ProLayout>
```

### 设置持久化

设置自动保存到浏览器的 `localStorage`，刷新页面后自动恢复：

```typescript
// 保存设置
localStorage.setItem('layout-settings', JSON.stringify(settings));

// 读取设置
const savedSettings = localStorage.getItem('layout-settings');
```

### 清除设置

如果需要恢复默认设置，可以在浏览器控制台执行：

```javascript
localStorage.removeItem('layout-settings');
location.reload();
```

## 自定义配置

### 修改默认设置

编辑 `src/components/Layout/MainLayout.tsx`：

```typescript
const [settings, setSettings] = useState({
  layout: 'side',        // 布局模式
  contentWidth: 'Fluid', // 内容宽度
  theme: 'light',        // 主题
  splitMenus: false,     // 分组菜单
  fixedHeader: false,    // 固定顶栏
  fixSiderbar: false,    // 固定侧边栏
});
```

### 添加更多主题色

在 `settingsProps.colorList` 中添加：

```typescript
colorList: [
  {
    key: 'custom',
    color: '#FF0000', // 自定义颜色
  },
  // ...其他颜色
]
```

### 修改菜单项

编辑 `menuItems` 数组：

```typescript
const menuItems = [
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
  // 添加新菜单项...
];
```

## API 参考

### ProLayout Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| title | string | - | 系统标题 |
| logo | boolean/ReactNode | false | Logo |
| route | RouteProps | - | 菜单路由配置 |
| location | Location | - | 当前路由位置 |
| theme | 'light' \| 'dark' | 'light' | 主题模式 |
| layout | 'side' \| 'top' \| 'mix' | 'side' | 布局模式 |
| contentWidth | 'Fluid' \| 'Fixed' | 'Fluid' | 内容宽度 |
| fixedHeader | boolean | false | 固定顶栏 |
| fixSiderbar | boolean | false | 固定侧边栏 |
| splitMenus | boolean | false | 分组菜单 |
| onSettingChange | (settings) => void | - | 设置变更回调 |
| settingsProps | SettingsProps | - | 设置面板配置 |

### SettingsProps

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| theme | 'light' \| 'dark' | - | 当前主题 |
| prefixCls | string | 'ant-pro' | 样式前缀 |
| settings | object | - | 当前设置 |
| onSettingChange | (settings) => void | - | 设置变更回调 |
| hideSettings | boolean | false | 隐藏设置按钮 |
| hideCopyButton | boolean | false | 隐藏复制按钮 |
| colorList | ColorItem[] | - | 主题色列表 |

## 常见问题

### Q: 设置没有生效？
A: 检查浏览器控制台是否有错误，清除缓存后重试。

### Q: 如何恢复默认设置？
A: 在控制台执行 `localStorage.removeItem('layout-settings'); location.reload();`

### Q: 主题色没有变化？
A: 确保已清除浏览器缓存，重新加载页面。

### Q: 暗色主题下某些组件样式不对？
A: 检查组件是否使用了硬编码的颜色值，应该使用 Ant Design 的主题变量。

## 最佳实践

1. **保持默认设置简洁**：默认配置应该适合大多数用户
2. **尊重用户选择**：用户修改的设置应该持久化
3. **提供预览**：修改设置时提供实时预览
4. **移动端适配**：在小屏幕上自动调整布局
5. **性能优化**：避免频繁写入 localStorage

## 相关资源

- [Ant Design Pro Components 文档](https://procomponents.ant.design/)
- [ProLayout 文档](https://procomponents.ant.design/components/layout)
- [Ant Design 主题定制](https://ant.design/docs/react/customize-theme-cn)
