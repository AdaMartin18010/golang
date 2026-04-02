# Go Workspaces (多模块工作区)

> **分类**: 开源技术堆栈  
> **标签**: #workspace #gomodules #go1.18

---

## 工作区配置

### go.work

```go
// go.work
go 1.22

use (
    ./module-a
    ./module-b
    ./module-c
)

replace (
    example.com/module-a => ./module-a
    example.com/module-b => ./module-b
)
```

---

## 使用场景

### 场景1: 多模块开发

```
project/
├── go.work
├── api/                    # 模块1: API 定义
│   ├── go.mod
│   └── proto/
├── server/                 # 模块2: 服务端
│   ├── go.mod
│   └── main.go
└── client/                 # 模块3: 客户端
    ├── go.mod
    └── main.go
```

```bash
# 初始化工作区
go work init ./api ./server ./client

# 现在可以在工作区内跨模块引用
cd server
go build  # 自动使用本地 api 模块
```

---

## 常用命令

```bash
# 初始化工作区
go work init [moddirs...]

# 添加模块到工作区
go work use ./new-module

# 从工作区移除模块
go work edit -dropuse=./old-module

# 查看工作区状态
cat go.work
```

---

## 与 replace 对比

| 特性 | go.work | go.mod replace |
|------|---------|----------------|
| 范围 | 本地开发 | 提交到仓库 |
| 团队协作 | 本地配置 | 共享配置 |
| CI/CD | 不使用 | 使用 |
| 适用场景 | 本地多模块开发 | 临时替换依赖 |

---

## CI/CD 处理

```bash
# CI 中忽略工作区，使用正式发布版本
rm go.work go.work.sum
go mod download
go build ./...
```

---

## 版本控制

```gitignore
# .gitignore
go.work
go.work.sum
```

工作区文件通常是本地配置，不应提交到版本控制。

---

## 最佳实践

1. **本地开发**使用 go.work
2. **CI/CD**删除 go.work，使用实际版本
3. **模块发布**后更新依赖
4. **文档**说明工作区使用方法
