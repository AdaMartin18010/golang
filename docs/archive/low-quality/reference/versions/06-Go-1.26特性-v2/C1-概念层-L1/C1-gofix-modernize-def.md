# go fix现代化 (go fix Modernization)

> **文档层级**: C1-概念层 (Concept Layer L1)
> **文档类型**: 概念定义 (Concept Definition)
> **工具**: `go fix`
> **最后更新**: 2026-03-06

---

## 一、概念定义

### 1.1 功能概述

```
go fix: 自动代码现代化工具

Go 1.26增强:
  - 自动应用新的语言特性
  - 更新废弃的API调用
  - 应用最佳实践模式

工作原理:
  1. 解析Go源代码
  2. 识别可现代化模式
  3. 应用转换规则
  4. 输出更新后的代码
```

### 1.2 支持的现代化规则

| 规则 | 说明 | 示例 |
|------|------|------|
| `newexpr` | 转换为new(expr)语法 | `&[]int{42}[0]` → `new(42)` |
| `errorsas` | 使用errors.AsType | `errors.As(err, &e)` → `errors.AsType[Error](err)` |
| `gofmt` | 格式优化 | 代码格式化 |

---

## 二、使用方法

### 2.1 基本命令

```bash
# 自动现代化当前目录代码
go fix ./...

# 查看将要做的变更（不实际修改）
go fix -diff ./...

# 现代化特定文件
go fix main.go
```

### 2.2 示例转换

```go
// 转换前 (Go 1.25风格)
config := Config{
    Timeout: &[]int{30}[0],
}

var myErr *MyError
if errors.As(err, &myErr) {
    // ...
}

// 转换后 (Go 1.26风格)
config := Config{
    Timeout: new(30),  // ✨ 更简洁
}

if myErr, ok := errors.AsType[*MyError](err); ok {
    // ✨ 类型安全
}
```

---

## 三、最佳实践

```bash
# 升级流程
1. 更新go.mod到1.26
echo 'go 1.26' > go.mod

2. 运行go fix
go fix ./...

3. 运行测试
go test ./...

4. 代码审查
git diff  # 检查变更
```

---

**概念分类**: 工具链 - 代码现代化
**Go版本**: 1.26+
**工具**: `go fix`
