# Go标准库

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go标准库](#go标准库)
  - [📋 目录](#-目录)
  - [📚 核心内容](#-核心内容)
    - [核心包](#核心包)
    - [I/O与文件](#io与文件)
    - [网络与并发](#网络与并发)
    - [工具与测试](#工具与测试)
  - [🚀 快速参考](#-快速参考)
  - [📖 系统文档](#-系统文档)

---

## 📚 核心内容

### 核心包

1. **[核心包概览](./01-核心包概览.md)** ⭐⭐⭐⭐⭐
   - fmt, io, os, time
   - strings, strconv
   - errors, log

### I/O与文件

- bufio, io, os
- path/filepath
- encoding (json, xml, csv)

### 网络与并发

- net/http, net
- Context, sync
- time

### 工具与测试

- testing, flag
- sort, math
- crypto, compress

---

## 🚀 快速参考

```go
// HTTP
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

// JSON
json.Marshal(data)
json.Unmarshal(data, &result)

// 文件
os.ReadFile("file.txt")
os.WriteFile("file.txt", data, 0644)
```

---

## 📖 系统文档

### 📚 核心系统文档

- **[知识图谱](./00-知识图谱.md)** - 标准库知识体系
- **[对比矩阵](./00-对比矩阵.md)** - 标准库包对比
- **[概念定义体系](./00-概念定义体系.md)** - 核心概念定义

### 📖 文档列表

1. **[核心包概览](./01-核心包概览.md)** - 标准库核心包详解

### 🎯 学习路径

**第一周：基础包**
- fmt, io, os → 基础I/O操作
- strings, strconv → 字符串处理
- time → 时间处理

**第二周：并发与网络**
- Context, sync → 并发控制
- net/http → HTTP编程
- encoding/json → JSON处理

**第三周：高级包**
- testing → 测试
- crypto → 加密
- database/sql → 数据库

### 🔗 相关资源

**相关主题**:
- [并发编程](../concurrency/) - 并发原语详解
- [Web开发](../../development/web/) - HTTP服务开发
- [数据库](../../development/database/) - 数据库操作

**外部资源**:
- [Go标准库文档](https://pkg.go.dev/std) - 官方文档
- [Go by Example](https://gobyexample.com/) - 标准库示例

---

**上次更新**: 2025-12-03
**维护者**: Go Framework Team
