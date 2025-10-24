# Go Formal Verification - Web UI Frontend

React + TypeScript + Vite前端应用，为Go形式化验证框架提供可视化界面。

## 🎯 功能特性

### 已实现

- ✅ React 18 + TypeScript
- ✅ Vite构建工具
- ✅ Tailwind CSS样式框架
- ✅ React Router路由管理
- ✅ Axios API客户端
- ✅ WebSocket实时通信
- ✅ 响应式布局
- ✅ 4个主要页面（Dashboard, Analysis, Patterns, Projects）

### 待实现

- [ ] D3.js CFG可视化
- [ ] 并发分析仪表板
- [ ] 交互式模式生成器
- [ ] 项目管理界面
- [ ] 实时分析进度
- [ ] 代码编辑器集成

## 🚀 快速开始

### 前置要求

- Node.js 18+
- npm 或 yarn

### 安装依赖

```bash
cd web-ui/frontend
npm install
```

### 开发模式

```bash
npm run dev
```

应用将在 `http://localhost:5173` 启动

### 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist/` 目录

### 预览生产构建

```bash
npm run preview
```

## 📁 项目结构

```text
frontend/
├── src/
│   ├── components/         # React组件
│   │   └── Layout.tsx      # 主布局组件
│   ├── pages/              # 页面组件
│   │   ├── Dashboard.tsx   # 仪表板
│   │   ├── Analysis.tsx    # 分析页面
│   │   ├── Patterns.tsx    # 模式页面
│   │   └── Projects.tsx    # 项目页面
│   ├── utils/              # 工具函数
│   │   ├── api.ts          # API客户端
│   │   └── websocket.ts    # WebSocket客户端
│   ├── App.tsx             # 根组件
│   ├── main.tsx            # 入口文件
│   └── index.css           # 全局样式
├── public/                 # 静态资源
├── index.html              # HTML模板
├── vite.config.ts          # Vite配置
├── tailwind.config.js      # Tailwind配置
├── tsconfig.json           # TypeScript配置
└── package.json            # 项目配置
```

## 🔌 API集成

### 后端连接

前端通过Vite代理连接后端API：

```typescript
// vite.config.ts
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
  },
  '/ws': {
    target: 'ws://localhost:8080',
    ws: true,
  },
}
```

### API客户端使用

```typescript
import { analysisAPI, patternsAPI, projectsAPI } from '@/utils/api'

// 分析代码
const result = await analysisAPI.analyzeCFG(code)

// 获取模式列表
const patterns = await patternsAPI.list()

// 创建项目
const project = await projectsAPI.create(name, description, path)
```

### WebSocket使用

```typescript
import { wsClient } from '@/utils/websocket'

// 监听消息
wsClient.on('connected', (data) => {
  console.log('Connected:', data)
})

wsClient.on('progress', (data) => {
  console.log('Progress:', data.progress)
})

// 发送消息
wsClient.send({ type: 'analyze', data: { code } })
```

## 🎨 样式系统

### Tailwind CSS

项目使用Tailwind CSS进行样式管理：

```tsx
// 主题颜色
<div className="bg-primary-600 text-white">
  Primary Button
</div>

// 响应式设计
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
  Responsive Grid
</div>
```

### 自定义主题

在 `tailwind.config.js` 中配置：

```javascript
theme: {
  extend: {
    colors: {
      primary: {
        500: '#0ea5e9',
        600: '#0284c7',
        700: '#0369a1',
      },
    },
  },
}
```

## 📊 页面说明

### Dashboard (仪表板)

- 项目统计
- 后端健康状态
- 快速操作入口
- 最近活动

### Analysis (分析页面)

- CFG可视化（待实现）
- 并发分析结果（待实现）
- 类型检查（待实现）

### Patterns (模式页面)

- 30+并发模式浏览（待实现）
- 交互式代码生成（待实现）
- CSP定义展示（待实现）

### Projects (项目页面)

- 项目列表（待实现）
- 项目管理（待实现）
- 分析历史（待实现）

## 🔧 开发指南

### 添加新页面

1. 在 `src/pages/` 创建新组件
2. 在 `src/App.tsx` 添加路由
3. 在 `src/components/Layout.tsx` 添加导航链接

### 添加新API

在 `src/utils/api.ts` 添加新的API方法：

```typescript
export const newAPI = {
  getData: () => apiClient.get('/new-endpoint'),
  postData: (data: any) => apiClient.post('/new-endpoint', data),
}
```

### 状态管理

当前使用React Hook进行本地状态管理。对于全局状态，可以考虑：

- Zustand (已安装)
- React Context
- Redux Toolkit

## 🧪 测试

```bash
# 运行测试
npm run test

# 运行测试并生成覆盖率报告
npm run test:coverage
```

## 📝 代码规范

### ESLint

```bash
# 检查代码
npm run lint

# 自动修复
npm run lint:fix
```

### TypeScript

项目使用严格的TypeScript配置：

- `strict: true`
- `noUnusedLocals: true`
- `noUnusedParameters: true`
- `noFallthroughCasesInSwitch: true`

## 🚧 开发状态

**当前版本**: v0.1.0 (Alpha)

### 完成度

- 基础架构: 100% ✅
- 页面框架: 100% ✅
- API集成: 100% ✅
- WebSocket: 100% ✅
- UI实现: 20% 🔄

### 下一步

1. 实现CFG可视化组件（D3.js）
2. 实现并发分析仪表板
3. 实现模式生成器UI
4. 完善项目管理功能

## 📞 联系方式

- Issues: <https://github.com/your-org/go-formal-verification/issues>
- Documentation: <https://your-org.github.io/go-formal-verification>

---

**Go Formal Verification Framework - Web UI Frontend**  
*从CLI到可视化，让形式化验证触手可及！* 🚀
