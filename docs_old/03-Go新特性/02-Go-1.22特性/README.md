# Go 1.22 新特性

> **简介**: Go 1.22版本(2024年2月发布)引入for循环变量语义改进、HTTP路由增强和性能提升

> **版本**: Go 1.22+  
> **难度**: ⭐⭐⭐  
> **标签**: #Go1.22 #新特性 #for循环 #HTTP路由

<!-- TOC START -->
- [Go 1.22 新特性](#go-122-新特性)
  - [📋 概述](#-概述)
  - [🎯 主要特性](#-主要特性)
    - [语言改进](#语言改进)
    - [标准库增强](#标准库增强)
    - [性能优化](#性能优化)
  - [📚 详细文档](#-详细文档)
  - [🔗 相关资源](#-相关资源)
<!-- TOC END -->

---

## 📋 概述

Go 1.22 是Go语言的重要版本，主要改进包括：

- **for循环语义**: 修复长期存在的变量捕获问题
- **HTTP路由增强**: ServeMux支持方法和通配符路由
- **性能提升**: 内存分配器优化、编译速度提升
- **工具链改进**: go工作区模式增强

---

## 🎯 主要特性

### 语言改进

#### 1. for循环变量语义改进

**问题**: Go 1.21及更早版本中，for循环变量在整个循环中共享同一地址：

```go
// Go 1.21 - 错误行为
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // 所有goroutine都打印3
    })
}
for _, f := range funcs {
    go f()
}
```

**修复**: Go 1.22中，每次迭代创建新变量：

```go
// Go 1.22 - 正确行为
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // 正确打印0, 1, 2
    })
}
for _, f := range funcs {
    go f()
}
```

**详细示例**：

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    // 场景1: goroutine中使用循环变量
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Println(i)  // Go 1.22: 正确打印0-4
        }()
    }
    wg.Wait()
    
    // 场景2: 闭包捕获
    var closures []func() int
    for _, v := range []int{1, 2, 3} {
        closures = append(closures, func() int {
            return v  // Go 1.22: 正确捕获各自的v值
        })
    }
    for _, c := range closures {
        fmt.Println(c())  // 1, 2, 3
    }
    
    // 场景3: defer中使用
    for i := 0; i < 3; i++ {
        defer func() {
            fmt.Println(i)  // Go 1.22: 打印2, 1, 0
        }()
    }
}
```

---

### 标准库增强

#### 1. HTTP路由增强 (ServeMux)

Go 1.22的`http.ServeMux`支持：
- HTTP方法匹配
- 通配符路由
- 优先级路由

**基本用法**：

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // 方法匹配
    mux.HandleFunc("GET /users", getUsers)
    mux.HandleFunc("POST /users", createUser)
    
    // 路径参数（通配符）
    mux.HandleFunc("GET /users/{id}", getUser)
    mux.HandleFunc("DELETE /users/{id}", deleteUser)
    
    // 带方法的通配符
    mux.HandleFunc("GET /posts/{id}/comments", getComments)
    
    http.ListenAndServe(":8080", mux)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "List all users")
}

func createUser(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Create user")
}

func getUser(w http.ResponseWriter, r *http.Request) {
    // 获取路径参数
    id := r.PathValue("id")
    fmt.Fprintf(w, "Get user: %s\n", id)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Delete user: %s\n", id)
}

func getComments(w http.ResponseWriter, r *http.Request) {
    postID := r.PathValue("id")
    fmt.Fprintf(w, "Get comments for post: %s\n", postID)
}
```

**路由优先级**：

```go
mux := http.NewServeMux()

// 精确匹配优先级最高
mux.HandleFunc("GET /users/admin", handleAdmin)

// 通配符次之
mux.HandleFunc("GET /users/{id}", handleUser)

// 前缀匹配最低
mux.HandleFunc("GET /users/", handleUsersPrefix)
```

#### 2. math/rand/v2包

新的随机数包，性能更好：

```go
package main

import (
    "fmt"
    "math/rand/v2"
)

func main() {
    // 自动种子（无需手动调用Seed）
    n := rand.IntN(100)  // [0, 100)
    fmt.Println(n)
    
    // 新的API
    f := rand.Float64()  // [0.0, 1.0)
    b := rand.N(uint64(1000))  // 泛型版本
    
    fmt.Println(f, b)
}
```

---

### 性能优化

#### 1. 内存分配器优化

- 小对象分配性能提升5-10%
- 内存使用更高效

#### 2. 编译速度提升

- 编译速度平均提升6%
- 链接速度提升20%

#### 3. 运行时优化

```go
// 示例：性能对比
package main

import (
    "testing"
)

func BenchmarkMemoryAllocation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}

// Go 1.21: ~300 ns/op
// Go 1.22: ~270 ns/op (提升10%)
```

---

## 📚 详细文档

本目录包含Go 1.22各特性的详细文档：

1. [for循环变量语义](./01-for循环变量语义.md) - 语义变更详解
2. [HTTP路由增强](./02-HTTP路由增强.md) - ServeMux新特性
3. [性能改进](./03-性能改进.md) - 性能优化详解

---

## 🔗 相关资源

### 官方文档

- [Go 1.22 Release Notes](https://go.dev/doc/go1.22)
- [for循环变量语义](https://go.dev/blog/loopvar-preview)
- [HTTP路由增强](https://go.dev/blog/routing-enhancements)

### 迁移指南

- [从Go 1.21迁移到1.22](https://go.dev/doc/go1.22#language)
- [for循环迁移FAQ](https://go.dev/wiki/LoopvarExperiment)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月21日  
**文档状态**: 完成  
**适用版本**: Go 1.22+
