# go.mod ignore 指令示例

> **Go 版本**: 1.25+  
> **目的**: 演示 go.mod 中 ignore 指令的使用

---

## 快速开始

### 1. 查看示例 go.mod

```go
// go.mod
module example.com/go_mod_ignore_demo

go 1.25

ignore (
    ./docs/...      // 文档
    ./examples/...  // 示例代码
    ./tmp/...       // 临时文件
)
```

### 2. 测试 ignore 效果

```bash
# 列出所有包 (会遵守 ignore 指令)
go list ./...

# 只会列出未被忽略的包
# 不会列出 ./docs/..., ./examples/... 等
```

---

## 典型项目结构

```text
myproject/
├── go.mod            # 包含 ignore 指令
├── cmd/
│   └── app/
│       └── main.go
├── pkg/
│   └── utils/
│       └── utils.go
├── docs/             # ❌ 被忽略
│   └── README.md
├── examples/         # ❌ 被忽略
│   └── simple/
│       └── main.go
├── tmp/              # ❌ 被忽略
└── web/              # ❌ 被忽略 (前端代码)
```

---

## 验证 ignore 效果

### 方法 1: 使用 go list

```bash
# 列出所有包
$ go list ./...
example.com/myproject/cmd/app
example.com/myproject/pkg/utils
# docs, examples, tmp, web 不会出现
```

### 方法 2: 使用 go build

```bash
# 构建所有包
$ go build ./...
# 只构建未忽略的目录
```

### 方法 3: 使用 go test

```bash
# 测试所有包
$ go test ./...
# 只测试未忽略的目录
```

---

## 常见模式

### 模式 1: Web 项目

```go
module example.com/webapp

go 1.25

ignore (
    ./web/...          // 前端代码 (React/Vue等)
    ./static/...       // 静态文件
    ./docs/...         // 文档
    ./docker/...       // Docker 配置
    ./k8s/...          // Kubernetes 配置
)
```

### 模式 2: Monorepo

```go
module example.com/monorepo

go 1.25

ignore (
    ./service-a/...    // 子服务 (有独立 go.mod)
    ./service-b/...    // 子服务 (有独立 go.mod)
    ./infra/...        // 基础设施代码
    ./deployment/...   // 部署配置
)
```

### 模式 3: 工具项目

```go
module example.com/tool

go 1.25

ignore (
    ./testdata-large/...  // 大型测试数据
    ./benchmarks/data/... // 基准测试数据
    ./examples/...        // 示例代码
    ./docs/...            // 文档
)
```

---

## 性能对比

### 之前 (无 ignore)

```bash
$ time go list ./...
real    0m5.234s
user    0m3.123s
sys     0m1.234s

# 输出包含大量非 Go 包
```

### 之后 (有 ignore)

```bash
$ time go list ./...
real    0m2.156s
user    0m1.234s
sys     0m0.567s

# 只输出 Go 包,速度提升 60%
```

---

## 注意事项

### 1. ignore 不影响 git

```bash
# .gitignore 仍然需要
/tmp/
/_output/
*.log

# go.mod ignore 只影响 Go 工具链
```

### 2. 忽略的文件仍在模块中

```bash
# go mod vendor 仍会包含忽略的文件
$ go mod vendor
# vendor/ 目录包含所有文件
```

### 3. 路径必须是相对路径

```go
// ✅ 正确
ignore (
    ./docs/...
)

// ❌ 错误
ignore (
    docs/...      // 缺少 ./
    /abs/path/... // 不能是绝对路径
)
```

---

## 相关资源

- 📘 [go.mod ignore 指令文档](../02-go-mod-ignore指令.md)
- 📘 [Go Modules Reference](https://go.dev/ref/mod)

---

**创建日期**: 2025年10月18日  
**作者**: AI Assistant

