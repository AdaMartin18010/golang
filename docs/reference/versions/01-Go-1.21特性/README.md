# Go 1.21 新特性

> **简介**: Go 1.21版本(2023年8月发布)引入泛型改进、性能优化和标准库重大更新
> **版本**: Go 1.21+  
> **难度**: ⭐⭐⭐  
> **标签**: #Go1.21 #新特性 #泛型 #性能优化

<!-- TOC START -->
- [Go 1.21 新特性](#go-121-新特性)
  - [📋 概述](#-概述)
  - [🎯 主要特性](#-主要特性)
    - [语言增强](#语言增强)
      - [1. 内置泛型函数](#1-内置泛型函数)
      - [2. 泛型类型推断改进](#2-泛型类型推断改进)
    - [标准库更新](#标准库更新)
      - [1. log/slog - 结构化日志](#1-logslog---结构化日志)
      - [2. slices包增强](#2-slices包增强)
      - [3. maps包增强](#3-maps包增强)
    - [工具链改进](#工具链改进)
      - [1. PGO (Profile-Guided Optimization)](#1-pgo-profile-guided-optimization)
      - [2. 向后兼容性](#2-向后兼容性)
      - [3. go工具改进](#3-go工具改进)
  - [📚 详细文档](#-详细文档)
  - [🔗 相关资源](#-相关资源)
    - [官方文档](#官方文档)
    - [社区资源](#社区资源)
<!-- TOC END -->

---

## 📋 概述

Go 1.21 是Go语言的重要里程碑版本，主要改进包括：

- **泛型增强**: min/max内置函数、clear函数
- **性能优化**: PGO（Profile-Guided Optimization）正式版
- **标准库**: 新增log/slog结构化日志包
- **工具链**: 向后兼容性改进

---

## 🎯 主要特性

### 语言增强

#### 1. 内置泛型函数

Go 1.21 新增三个内置泛型函数：

```go
// min 返回最小值
func min[T cmp.Ordered](x, y T) T

// max 返回最大值  
func max[T cmp.Ordered](x, y T) T

// clear 清空map或slice
func clear[T ~[]Type | ~map[Type]Type1](t T)
```

**示例**：

```go
package main

import "fmt"

func main() {
    // min/max函数
    a := min(10, 20)     // 10
    b := max(10, 20)     // 20
    c := min(1.5, 2.3)   // 1.5
    
    fmt.Println(a, b, c)
    
    // clear函数
    m := map[string]int{"a": 1, "b": 2}
    clear(m)  // m现在为空
    fmt.Println(len(m))  // 0
    
    s := []int{1, 2, 3, 4, 5}
    clear(s)  // 所有元素置为零值
    fmt.Println(s)  // [0 0 0 0 0]
}
```

#### 2. 泛型类型推断改进

更智能的类型推断，减少显式类型参数：

```go
package main

func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

func main() {
    nums := []int{1, 2, 3}
    
    // Go 1.21: 自动推断类型
    strs := Map(nums, func(n int) string {
        return fmt.Sprintf("%d", n)
    })
    
    fmt.Println(strs)  // ["1" "2" "3"]
}
```

---

### 标准库更新

#### 1. log/slog - 结构化日志

Go 1.21 新增官方结构化日志包：

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // 默认logger
    slog.Info("application started", "version", "1.0.0", "port", 8080)
    
    // JSON格式
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger.Info("user login",
        slog.String("user", "alice"),
        slog.Int("user_id", 123),
    )
    
    // 分组
    logger.Info("request processed",
        slog.Group("http",
            slog.String("method", "GET"),
            slog.String("path", "/api/users"),
            slog.Int("status", 200),
        ),
    )
}
```

#### 2. slices包增强

新增实用切片操作：

```go
package main

import (
    "fmt"
    "slices"
)

func main() {
    s := []int{3, 1, 4, 1, 5, 9, 2, 6}
    
    // 排序
    slices.Sort(s)
    fmt.Println(s)  // [1 1 2 3 4 5 6 9]
    
    // 查找
    found := slices.Contains(s, 5)  // true
    idx := slices.Index(s, 5)       // 5
    
    // 去重
    unique := slices.Compact(s)
    fmt.Println(unique)  // [1 2 3 4 5 6 9]
}
```

#### 3. maps包增强

新增map操作函数：

```go
package main

import (
    "fmt"
    "maps"
)

func main() {
    m1 := map[string]int{"a": 1, "b": 2}
    m2 := map[string]int{"b": 3, "c": 4}
    
    // 克隆
    m3 := maps.Clone(m1)
    
    // 合并（后者覆盖前者）
    maps.Copy(m3, m2)
    fmt.Println(m3)  // {"a": 1, "b": 3, "c": 4}
    
    // 比较
    equal := maps.Equal(m1, m2)  // false
    fmt.Println(equal)
}
```

---

### 工具链改进

#### 1. PGO (Profile-Guided Optimization)

正式版性能优化：

```bash
# 生成profile
go test -cpuprofile=cpu.prof

# 使用PGO构建
go build -pgo=cpu.prof
```

性能提升：2-14%

#### 2. 向后兼容性

go.mod中指定Go版本：

```go
module example.com/myapp

go 1.21

require (
    // 依赖...
)
```

#### 3. go工具改进

```bash
# 查看所有可用版本
go list -m -versions golang.org/x/tools

# 工作区模式改进
go work use ./moduleA ./moduleB
```

---

## 📚 详细文档

本目录包含Go 1.21各特性的详细文档：

1. [内置函数详解](./01-内置函数详解.md) - min/max/clear函数完整指南
2. [PGO使用指南](./02-PGO使用指南.md) - Profile-Guided Optimization实战
3. [标准库增强详解](./03-标准库增强详解.md) - log/slog、slices、maps、cmp包详解

---

## 🔗 相关资源

### 官方文档

- [Go 1.21 Release Notes](https://go.dev/doc/go1.21)
- [log/slog包文档](https://pkg.go.dev/log/slog)
- [PGO用户指南](https://go.dev/doc/pgo)

### 社区资源

- [Go 1.21新特性详解](https://go.dev/blog/go1.21)
- [结构化日志最佳实践](https://go.dev/blog/slog)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月21日  
**文档状态**: 完成  
**适用版本**: Go 1.21+
