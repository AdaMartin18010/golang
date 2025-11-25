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
