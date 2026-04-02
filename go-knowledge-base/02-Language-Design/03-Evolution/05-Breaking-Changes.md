# 破坏性变更 (Breaking Changes)

> **分类**: 语言设计

---

## Go 1 兼容性承诺

> "It is intended that programs written to the Go 1 specification will continue to compile and run correctly, unchanged, over the lifetime of that specification."

---

## 兼容规则

### 1. 源代码兼容

```go
// 旧代码继续编译
package main

func main() {
    // Go 1.0 代码在 Go 1.26 仍可编译
}
```

### 2. API 兼容

标准库 API 保持稳定。

### 3. 例外情况

- 安全修复
- 规范错误修复
- 操作系统/架构弃用

---

## 已知破坏性变更

### 1. Go 1.22 循环变量

```go
// 1.21: 循环变量共享
// 1.22: 每次迭代新变量

// 旧代码依赖共享行为的可能需要调整
for i, v := range items {
    // 闭包捕获行为改变
}
```

**迁移**: Go 1.22 自动处理，但需注意旧代码假设。

### 2. 包路径变更

```go
// golang.org/x 包可能迁移
// 旧: golang.org/x/net/context
// 新: context (标准库)
```

### 3. 运行时行为

```go
// 垃圾回收行为可能改变
// 但不影响正确性
```

---

## 管理依赖

### go.mod

```go
module example.com/app

go 1.22  // 指定 Go 版本

require (
    // 依赖版本锁定
)
```

### 版本策略

1. **使用 go.mod 指定版本**
2. **定期更新依赖**
3. **测试升级**

---

## 未来风险

| 领域 | 潜在变更 |
|------|----------|
| 错误处理 | try/catch 语法讨论中 |
| 泛型 | 约束语法演进 |
| 并发 | 结构化并发 API |

---

## 最佳实践

1. **锁定 Go 版本**: go.mod
2. **持续集成**: 测试多个版本
3. **关注发布说明**: go.dev/doc/devel/release
