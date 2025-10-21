# Go 1.23 新特性

> **简介**: Go 1.23版本(2024年8月发布)引入迭代器预览、工具链增强和标准库改进

> **版本**: Go 1.23+  
> **难度**: ⭐⭐⭐⭐  
> **标签**: #Go1.23 #新特性 #迭代器 #工具链

<!-- TOC START -->
- [Go 1.23 新特性](#go-123-新特性)
  - [📋 概述](#-概述)
  - [🎯 主要特性](#-主要特性)
    - [语言预览](#语言预览)
    - [标准库增强](#标准库增强)
    - [工具链改进](#工具链改进)
  - [📚 详细文档](#-详细文档)
  - [🔗 相关资源](#-相关资源)
<!-- TOC END -->

---

## 📋 概述

Go 1.23 是Go语言向2.0演进的重要版本，主要改进包括：

- **迭代器预览**: range-over-func实验性支持
- **工具链增强**: go命令和工作区改进
- **标准库更新**: 新API和性能优化
- **向后兼容**: 保持API稳定性

---

## 🎯 主要特性

### 语言预览

#### 1. Range-over-func (实验性)

Go 1.23引入迭代器实验性支持（需要GOEXPERIMENT=rangefunc）：

```go
package main

import (
    "fmt"
    "iter"
)

// 自定义迭代器
func Count(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := 0; i < n; i++ {
            if !yield(i) {
                return
            }
        }
    }
}

func main() {
    // 使用迭代器
    for i := range Count(5) {
        fmt.Println(i)  // 0, 1, 2, 3, 4
    }
}
```

**更复杂的示例**：

```go
package main

import (
    "fmt"
    "iter"
)

// 键值对迭代器
func MapIter[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
    return func(yield func(K, V) bool) {
        for k, v := range m {
            if !yield(k, v) {
                return
            }
        }
    }
}

// 过滤迭代器
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
    return func(yield func(T) bool) {
        for v := range seq {
            if pred(v) {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

func main() {
    m := map[string]int{"a": 1, "b": 2, "c": 3}
    
    // 遍历map
    for k, v := range MapIter(m) {
        fmt.Printf("%s: %d\n", k, v)
    }
    
    // 过滤
    numbers := func(yield func(int) bool) {
        for i := 1; i <= 10; i++ {
            if !yield(i) {
                return
            }
        }
    }
    
    for n := range Filter(numbers, func(n int) bool { return n%2 == 0 }) {
        fmt.Println(n)  // 2, 4, 6, 8, 10
    }
}
```

---

### 标准库增强

#### 1. slices包新增函数

```go
package main

import (
    "fmt"
    "slices"
)

func main() {
    s := []int{1, 2, 3, 4, 5}
    
    // 反转
    slices.Reverse(s)
    fmt.Println(s)  // [5, 4, 3, 2, 1]
    
    // 按块处理
    chunks := slices.Chunk([]int{1, 2, 3, 4, 5, 6, 7}, 3)
    for chunk := range chunks {
        fmt.Println(chunk)  // [1,2,3], [4,5,6], [7]
    }
    
    // 去重并保持顺序
    unique := slices.CompactFunc([]int{1, 2, 2, 3, 3, 3, 4}, func(a, b int) bool {
        return a == b
    })
    fmt.Println(unique)  // [1, 2, 3, 4]
}
```

#### 2. maps包增强

```go
package main

import (
    "fmt"
    "maps"
)

func main() {
    m1 := map[string]int{"a": 1, "b": 2}
    m2 := map[string]int{"c": 3, "d": 4}
    
    // 删除满足条件的键值对
    maps.DeleteFunc(m1, func(k string, v int) bool {
        return v > 1
    })
    fmt.Println(m1)  // {"a": 1}
    
    // 合并map
    result := maps.Clone(m1)
    maps.Copy(result, m2)
    fmt.Println(result)  // {"a": 1, "c": 3, "d": 4}
}
```

#### 3. testing/slogtest包

用于测试slog handler实现：

```go
package main

import (
    "log/slog"
    "testing"
    "testing/slogtest"
)

func TestCustomHandler(t *testing.T) {
    handler := &MyCustomHandler{}
    
    // 验证handler实现是否正确
    err := slogtest.TestHandler(handler, func() []map[string]any {
        return nil  // 返回预期的日志记录
    })
    
    if err != nil {
        t.Fatal(err)
    }
}
```

---

### 工具链改进

#### 1. go命令增强

```bash
# 新的go env -changed命令
go env -changed  # 只显示修改过的环境变量

# go test增强
go test -fullpath  # 显示完整路径
go test -json=short  # 简化JSON输出

# go build优化
go build -cover  # 构建支持覆盖率的二进制
```

#### 2. 工作区改进

```bash
# 自动同步工作区
go work sync

# 工作区use命令增强
go work use -r ./...  # 递归添加所有模块
```

#### 3. pprof增强

```bash
# 新的内存分析模式
go tool pprof -alloc_space  # 总分配空间
go tool pprof -alloc_objects  # 总分配对象数
```

---

## 📚 详细文档

本目录包含Go 1.23各特性的详细文档：

1. [迭代器预览](./01-迭代器预览.md) - range-over-func详解
2. [工具链增强](./02-工具链增强.md) - go命令新特性
3. [标准库更新](./03-标准库更新.md) - API变更详解

---

## 🔗 相关资源

### 官方文档

- [Go 1.23 Release Notes](https://go.dev/doc/go1.23)
- [迭代器提案](https://github.com/golang/go/issues/61897)
- [GOEXPERIMENT设置](https://go.dev/doc/godebug)

### 实验性特性

- [启用rangefunc实验](https://go.dev/wiki/RangefuncExperiment)
- [迭代器最佳实践](https://go.dev/blog/range-functions)

### 迁移指南

- [从Go 1.22迁移到1.23](https://go.dev/doc/go1.23#language)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月21日  
**文档状态**: 完成  
**适用版本**: Go 1.23+
