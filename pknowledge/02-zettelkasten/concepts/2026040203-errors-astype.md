---
id: 2026040203
title: errors.AsType - 泛型错误处理
date: 2026-04-02
tags: [error-handling, generics, go126, performance]
references: [go126-release-notes]
status: active
---

## 语法

```go
if valErr, ok := errors.AsType[*ValidationError](err); ok {
    // valErr 已确定类型
}
```

## 与 errors.As 对比

| 特性 | errors.As | errors.AsType |
|------|-----------|---------------|
| 语法 | `var e *MyError; errors.As(err, &e)` | `e, ok := errors.AsType[*MyError](err)` |
| 类型安全 | 运行时 | 编译时 |
| 性能 | 95.62ns (反射) | 30.26ns (泛型) |
| 简洁性 | 需预声明变量 | 内联使用 |

**性能提升: 68%**

## 实现原理

使用泛型单态化 (monomorphization) 替代反射。

```go
// 编译器为每个具体类型生成专用代码
func AsType_ValidationError(err error) (*ValidationError, bool)
```

## 最佳实践

### Do ✅

```go
// 新代码使用 AsType
if pathErr, ok := errors.AsType[*fs.PathError](err); ok {
    log.Printf("path: %s", pathErr.Path)
}
```

### Don't ❌

```go
// 避免使用旧方式
var pathErr *fs.PathError
if errors.As(err, &pathErr) {
    // ...
}
```

## 迁移策略

1. 新代码全部使用 AsType
2. 旧代码逐步迁移
3. 性能敏感路径优先

## 关联

- [[Go Error Handling Patterns]]
- [[Generics Implementation]]
- [[Performance Optimization]]

## 代码示例

见 [[Error Handling Examples]]
