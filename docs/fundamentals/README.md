# Go基础知识

Go语言基础知识，涵盖语言特性、并发编程、标准库和数据结构。

---

## 📋 目录结构

### 核心模块

1. **[语言特性](./language/README.md)** ⭐⭐⭐⭐⭐
   - 语法基础
   - 并发编程
   - 模块管理
   - Go 1.25.3完整特性

2. **[并发编程](./concurrency/README.md)** ⭐⭐⭐⭐⭐
   - Goroutine与Channel
   - 并发模式
   - Context应用
   - 并发最佳实践

3. **[标准库](./stdlib/README.md)** ⭐⭐⭐⭐⭐
   - 核心包概览
   - 常用API参考

4. **[数据结构](./data-structures/README.md)** ⭐⭐⭐⭐
   - 基础数据结构
   - 高级数据结构
   - 算法实现

---

## 🎯 学习路径

### 初学者 (1-2周)

```text
语法基础 → 数据类型 → 控制流 → 函数
```

### 进阶者 (2-3周)

```text
结构体与接口 → 错误处理 → 包管理 → 测试
```

### 专家 (4-6周)

```text
并发编程 → 标准库深入 → 性能优化 → 源码阅读
```

---

## 🚀 快速开始

### Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go 1.25.3!")
}
```

### 并发示例

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch <- "Hello from goroutine!"
    }()
    
    msg := <-ch
    fmt.Println(msg)
}
```

---

## 📖 系统文档

- **[知识图谱](./00-知识图谱.md)**: 基础知识体系全景图
- **[对比矩阵](./00-对比矩阵.md)**: 不同方案对比
- **[概念定义体系](./00-概念定义体系.md)**: 核心概念详解

---

## 🔗 相关资源

- [Go 1.25.3完整知识体系](./language/00-Go-1.25.3完整知识体系总览-2025.md)
- [形式化理论体系](./language/00-Go-1.25.3形式化理论体系/README.md)
- [核心机制解析](./language/00-Go-1.25.3核心机制完整解析/README.md)

---

## 📚 推荐阅读顺序

1. **语言基础** → 语法基础 → 类型系统 → 控制流
2. **并发编程** → Goroutine → Channel → Context
3. **标准库** → 核心包 → 常用API
4. **数据结构** → 基础结构 → 算法实现

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
