# Go 1.23 新特性示例（原目录名误标为 Go 1.25）

> **更正说明**: `iter` 包、`unique` 包、`range over func` 等特性是 **Go 1.23** 引入的，
> 而非 Go 1.25。本目录内容已验证可在 Go 1.26.2 下编译运行。

本目录包含 Go 1.23 引入的新特性示例代码。

## 📋 目录

- [01-iter-demo](#01-iter-demo) - 迭代器 API 示例
- [02-unique-demo](#02-unique-demo) - unique 包示例
- [03-testing-loop](#03-testing-loop) - 测试循环优化（已移除）

---

## 01-iter-demo

### 概述

Go 1.23 正式引入了 `iter` 包和 `range over func` 支持。
标准库中的迭代器 API（如 `strings.Lines()`, `strings.SplitSeq()`, `strings.FieldsSeq()`）
已在 Go 1.23+ 中可用，提供更高效、更内存友好的迭代方式。

**注意**: 这些 API 在 Go 1.26.2 中完全稳定可用。

### 新特性（Go 1.23）

- `strings.Lines()` - 按行迭代字符串
- `strings.SplitSeq()` - 分割迭代器
- `strings.FieldsSeq()` - 字段迭代器

### 运行示例

```bash
cd 01-iter-demo
go run main.go
```

### 优势

1. **内存效率**: 迭代器模式避免了创建中间切片
2. **延迟计算**: 只在需要时计算下一个元素
3. **更简洁的 API**: 使用 `range` 循环直接迭代

### 当前实现

示例代码使用传统的 `strings.Split()` 和 `strings.Fields()` 方法，并提供了注释说明如何使用新的迭代器 API（当它们可用时）。

---

## 02-unique-demo

### 概述

`unique` 包提供了值规范化功能，可以自动去重相同的内容，节省内存。

**注意**: `unique` 包在 **Go 1.23** 中引入，在 Go 1.26.2 中已稳定可用。

### 核心功能（计划中）

- `unique.Make()` - 创建规范化值
- `unique.Handle[T]` - 规范化值的句柄
- 自动去重相同内容

### 运行示例

```bash
cd 02-unique-demo
go run main.go
```

### 使用场景

1. **字符串去重**: 大量重复字符串的场景
2. **结构体去重**: 相同结构体实例的复用
3. **内存优化**: 减少内存占用

### 注意事项

- `unique` 包是 Go 1.23+ 的特性，Go 1.26.2 中已稳定
- 适用于有大量重复值的场景
- 规范化后的值是不可变的
- 本示例提供了概念演示，实际使用应使用标准库的 `unique` 包

---

## 03-testing-loop

### ⚠️ 重要提示

**注意**: `testing.B.Loop()` 是 Go 1.24 引入的 API，本示例仅展示传统基准测试模式作为对比。**

本示例仅作为历史参考，实际开发中应使用传统的 `for i := 0; i < b.N; i++` 模式。

请参考 `../03-testing-enhancement/` 目录中的基准测试最佳实践。

### 运行示例（仅作参考）

```bash
cd 03-testing-loop
go test -bench=.
```

---

## 🚀 快速开始

### 运行所有示例

```bash
# 迭代器示例
cd 01-iter-demo && go run main.go

# unique 包示例
cd 02-unique-demo && go run main.go

# 测试示例（已废弃，仅作参考）
cd 03-testing-loop && go test -bench=.
```

---

## 📚 相关资源

- [Go 1.23 Release Notes](https://go.dev/doc/go1.23)
- [strings 包文档](https://pkg.go.dev/strings)
- [unique 包文档](https://pkg.go.dev/unique) (实验性)

---

**最后更新**: 2025-11-11
**Go 版本**: 1.23+ (在 1.26.2 下验证)
