# Go 1.25特性

Go 1.25版本特性完整指南，涵盖容器感知、greenteagc、语言改进和标准库更新。

---

## 📋 目录

- [Go 1.25特性](#go-125特性)
  - [📋 目录](#-目录)
  - [🎯 核心特性](#-核心特性)
    - [1. 容器感知型GOMAXPROCS ⭐⭐⭐⭐⭐](#1-容器感知型gomaxprocs-)
    - [2. greenteagc实验性GC ⭐⭐⭐⭐⭐](#2-greenteagc实验性gc-)
    - [3. Nil Pointer Panic改进 ⭐⭐⭐⭐](#3-nil-pointer-panic改进-)
    - [4. DWARF v5调试支持 ⭐⭐⭐⭐](#4-dwarf-v5调试支持-)
    - [5. 移除"core types"概念 ⭐⭐⭐](#5-移除core-types概念-)
  - [📚 详细文档](#-详细文档)
  - [🔧 迁移指南](#-迁移指南)
    - [从Go 1.24升级](#从go-124升级)
  - [🎯 最佳实践](#-最佳实践)
    - [容器部署](#容器部署)
    - [Kubernetes配置](#kubernetes配置)
  - [🔗 相关资源](#-相关资源)

## 🎯 核心特性

### 1. 容器感知型GOMAXPROCS ⭐⭐⭐⭐⭐

**自动适配容器资源限制**:

```go
// Go 1.25会自动检测Linux cgroup限制
// 无需手动设置GOMAXPROCS

// 之前 (Go 1.24)
import "runtime"
runtime.GOMAXPROCS(runtime.NumCPU())

// 现在 (Go 1.25)
// 自动适配，无需手动设置
```

**场景**:

- Kubernetes Pod CPU限制
- Docker容器CPU限制
- 避免过度调度

**效果**:

- 性能提升5-15%
- 避免CPU争抢
- 更好的资源利用

---

### 2. greenteagc实验性GC ⭐⭐⭐⭐⭐

**新型垃圾回收器**:

```bash
# 启用greenteagc
GOEXPERIMENT=greenteagc go build

# 环境变量
export GOEXPERIMENT=greenteagc
```

**优势**:

- GC开销降低10-40%
- 更低的延迟
- 更少的STW时间

**适用场景**:

- 高并发服务
- 低延迟要求
- 大内存应用

---

### 3. Nil Pointer Panic改进 ⭐⭐⭐⭐

**更准确的错误提示**:

```go
// Go 1.24
panic: runtime error: invalid memory address or nil pointer dereference

// Go 1.25
panic: runtime error: nil pointer dereference at field User.Name
```

**改进**:

- 更精确的字段定位
- 更快的调试速度
- 更好的错误信息

---

### 4. DWARF v5调试支持 ⭐⭐⭐⭐

**调试信息升级**:

```bash
go build -gcflags="-dwarfv=5" -o myapp
```

**优势**:

- 二进制大小减少30%
- 更快的调试速度
- 更好的IDE支持

---

### 5. 移除"core types"概念 ⭐⭐⭐

**简化泛型规范**:

```go
// Go 1.24及之前需要理解"core types"
type Number interface {
    ~int | ~float64  // core type是什么？
}

// Go 1.25简化了规范
type Number interface {
    ~int | ~float64  // 直接理解为类型约束
}
```

**改进**:

- 更简单的泛型概念
- 更容易理解的规范
- 更少的学习曲线

---

## 📚 详细文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)
- [实践应用](./README.md)

---

## 🔧 迁移指南

### 从Go 1.24升级

**1. 容器环境优化**:

```bash
# 删除手动GOMAXPROCS设置
# runtime.GOMAXPROCS(n) // 不再需要

# Go 1.25自动处理
```

**2. 启用greenteagc**:

```bash
# Dockerfile
ENV GOEXPERIMENT=greenteagc
RUN go build -o myapp

# 或构建时
GOEXPERIMENT=greenteagc go build
```

**3. 利用改进的错误信息**:

```go
// 添加更好的错误处理
if user == nil {
    return errors.New("user is nil")
}
// Go 1.25的panic会自动显示更详细的信息
```

---

## 🎯 最佳实践

### 容器部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .

# 启用greenteagc
ENV GOEXPERIMENT=greenteagc
RUN go build -o myapp

FROM alpine:latest
COPY --from=builder /app/myapp .
# GOMAXPROCS自动适配，无需设置
CMD ["./myapp"]
```

### Kubernetes配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:go1.25
        resources:
          limits:
            cpu: "2"  # Go 1.25自动适配
          requests:
            cpu: "1"
```

---

## 🔗 相关资源

- [Go 1.25发布说明](https://go.dev/doc/go1.25)
- [容器感知详解](https://go.dev/blog/container-aware-gomaxprocs)
- [greenteagc介绍](https://go.dev/wiki/greenteagc)
- [版本对比](../00-版本对比与选择指南.md)

---

**发布时间**: 2025年8月
**最后更新**: 2025-10-29
**当前版本**: 1.25.3
