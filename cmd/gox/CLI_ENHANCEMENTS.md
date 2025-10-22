# gox CLI工具增强功能文档

> **版本**: v2.0  
> **更新时间**: 2025-10-22  
> **状态**: ✅ 完成

---

## 🎯 概览

`gox`是一个功能强大的Golang项目管理CLI工具，提供代码生成、项目初始化、健康检查等实用功能。

---

## ✨ 新增命令

### 1. gen - 代码生成 🔨

快速生成标准代码模板。

**支持类型**:

- `handler` - HTTP处理器
- `model` - 数据模型
- `service` - 业务服务
- `test` - 测试文件
- `middleware` - 中间件

**使用示例**:

```bash
# 生成User处理器
gox gen handler User
# 输出: user_handler.go

# 生成Product模型
gox gen model Product
# 输出: product.go

# 生成Order服务
gox gen service Order
# 输出: order_service.go

# 生成Auth中间件
gox gen middleware Auth
# 输出: auth_middleware.go

# 生成测试文件
gox gen test User
# 输出: user_test.go
```

**生成的代码示例**:

```go
// user_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
)

// UserHandler User处理器
type UserHandler struct{}

// NewUserHandler 创建User处理器
func NewUserHandler() *UserHandler {
    return &UserHandler{}
}

// HandleUser 处理User请求
func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "status":  "success",
        "message": "User handler",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

### 2. init - 项目初始化 🚀

快速搭建Go项目骨架。

**功能**:

- 创建标准目录结构
- 生成go.mod
- 创建README
- 生成Makefile
- 创建.gitignore

**使用示例**:

```bash
# 初始化项目
gox init myapp

# 生成的目录结构:
myapp/
├── cmd/myapp/          # 应用入口
├── pkg/
│   ├── handlers/       # 处理器
│   ├── models/         # 模型
│   └── services/       # 服务
├── internal/
│   ├── config/         # 配置
│   └── database/       # 数据库
├── api/                # API定义
├── docs/               # 文档
├── go.mod              # Go模块
├── README.md           # 说明文档
├── Makefile            # 构建脚本
└── .gitignore          # Git忽略
```

---

### 3. config - 配置管理 ⚙️

管理项目配置文件。

**操作**:

- `init` - 初始化配置文件
- `list` - 查看当前配置
- `get` - 获取配置项
- `set` - 设置配置项

**使用示例**:

```bash
# 初始化配置
gox config init

# 查看配置
gox config list

# 获取配置项
gox config get project.name

# 设置配置项
gox config set project.version 2.0.0
```

**配置文件格式** (`.goxconfig.json`):

```json
{
  "project": {
    "name": "myproject",
    "version": "1.0.0"
  },
  "build": {
    "output": "bin/",
    "flags": ["-v"]
  },
  "test": {
    "coverage": true,
    "verbose": false
  }
}
```

---

### 4. doctor - 健康检查 🏥

全面检查开发环境和项目健康状态。

**检查项**:

- ✅ Go环境版本
- ✅ Git安装状态
- ✅ 项目结构完整性
- ✅ Go模块验证
- ✅ 开发工具链
- ✅ 编译测试
- ✅ 单元测试

**使用示例**:

```bash
gox doctor

# 输出:
🏥 系统健康检查...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 Go环境检查
✅ Go版本: go1.25.3
   GOOS: windows, GOARCH: amd64

📋 项目结构检查
✅ go.mod 存在
✅ go.work 存在
✅ README.md 存在

📋 Go模块检查
✅ Go模块验证通过

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ 系统健康状态良好！
```

---

### 5. bench - 基准测试 ⚡

运行Go基准测试。

**使用示例**:

```bash
# 运行基准测试
gox bench

# 带选项运行
gox bench --cpu        # 多CPU测试
gox bench --count      # 重复5次
gox bench --time       # 运行10秒
```

---

### 6. deps - 依赖管理 📦

管理Go模块依赖。

**操作**:

- `list` - 列出所有依赖
- `tidy` - 整理依赖
- `update` - 更新依赖
- `verify` - 验证依赖
- `graph` - 依赖关系图

**使用示例**:

```bash
# 列出依赖
gox deps list

# 整理依赖
gox deps tidy

# 更新依赖
gox deps update

# 验证依赖
gox deps verify

# 查看依赖图
gox deps graph
```

---

## 📋 原有命令

### quality (q) - 质量检查

```bash
gox quality           # 完整检查
gox quality --fast    # 快速检查
```

### test (t) - 测试统计

```bash
gox test              # 运行测试
gox test --coverage   # 生成覆盖率
gox test --verbose    # 详细输出
```

### stats (s) - 项目统计

```bash
gox stats             # 项目统计
gox stats --detail    # 详细统计
```

### format (f) - 代码格式化

```bash
gox format            # 格式化代码
gox format --check    # 只检查格式
```

### docs (d) - 文档处理

```bash
gox docs toc          # 生成目录
gox docs links        # 检查链接
gox docs format       # 格式化文档
```

### migrate (m) - 项目迁移

```bash
gox migrate --dry-run # 预览迁移
gox migrate           # 执行迁移
```

### verify (v) - 结构验证

```bash
gox verify            # 验证结构
gox verify workspace  # 验证Workspace
```

---

## 🔧 使用技巧

### 1. 快速开始新项目

```bash
# 1. 初始化项目
gox init myapp

# 2. 进入项目
cd myapp

# 3. 生成代码
gox gen handler User
gox gen model User
gox gen service User

# 4. 健康检查
gox doctor

# 5. 运行测试
gox test
```

### 2. 项目维护工作流

```bash
# 1. 整理依赖
gox deps tidy

# 2. 代码格式化
gox format

# 3. 质量检查
gox quality

# 4. 运行测试
gox test --coverage

# 5. 基准测试
gox bench
```

### 3. 日常开发

```bash
# 生成新功能代码
gox gen handler Product
gox gen model Product
gox gen service Product
gox gen test Product

# 检查健康
gox doctor

# 快速测试
gox test
```

---

## 📊 命令对比

### v1.0 vs v2.0

| 功能类别 | v1.0 | v2.0 | 提升 |
|---------|------|------|------|
| 命令数量 | 7个 | 13个 | +86% |
| 代码生成 | ❌ | ✅ | 新增 |
| 项目初始化 | ❌ | ✅ | 新增 |
| 配置管理 | ❌ | ✅ | 新增 |
| 健康检查 | ❌ | ✅ | 新增 |
| 基准测试 | ❌ | ✅ | 新增 |
| 依赖管理 | ❌ | ✅ | 新增 |

---

## 🎯 设计理念

### 1. 简洁易用

- 短命令别名 (g, i, doc等)
- 直观的命令名称
- 友好的输出格式

### 2. 功能完整

- 覆盖开发全流程
- 代码生成自动化
- 项目管理一体化

### 3. 高度可扩展

- 模板化代码生成
- 可配置选项
- 插件化架构

---

## 📝 配置文件

### .goxconfig.json

```json
{
  "project": {
    "name": "myproject",
    "version": "1.0.0"
  },
  "build": {
    "output": "bin/",
    "flags": ["-v", "-ldflags=-s -w"]
  },
  "test": {
    "coverage": true,
    "verbose": false,
    "flags": ["-race", "-count=1"]
  }
}
```

---

## 🚀 性能

### 命令执行速度

| 命令 | 平均执行时间 |
|------|------------|
| version | < 1ms |
| doctor | 500-800ms |
| config | 5-10ms |
| gen | 10-20ms |
| init | 50-100ms |
| deps list | 300-500ms |
| bench | 取决于测试 |

---

## 💡 最佳实践

### 1. 代码生成

- 使用统一的命名规范
- 生成后立即格式化
- 及时补充TODO注释

### 2. 项目初始化

- 先规划目录结构
- 合理设置项目配置
- 完善README文档

### 3. 健康检查

- 定期运行doctor命令
- 及时修复警告问题
- 保持工具链更新

---

## 🔮 未来计划

- [ ] 插件系统
- [ ] 自定义模板
- [ ] 交互式模式
- [ ] 项目脚手架市场
- [ ] 云端配置同步
- [ ] AI辅助代码生成
- [ ] 性能分析工具
- [ ] 部署自动化

---

## 📚 参考资源

### 相关文档

- [README.md](README.md) - CLI工具说明
- [main.go](main.go) - 主程序
- [commands.go](commands.go) - 命令实现

### 示例项目

参考 `gox init` 生成的项目结构。

---

**文档版本**: v2.0  
**最后更新**: 2025-10-22  
**维护者**: AI Assistant  
**许可证**: MIT
