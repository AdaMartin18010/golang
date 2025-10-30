# Go标准库

Go标准库完整参考，涵盖核心包和常用API。

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
- context, sync
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

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3
