# API参考

Go API参考文档，涵盖标准库、第三方库和API设计最佳实践。

---


## 📋 目录


- [📚 文档导航](#-文档导航)
  - [核心内容](#核心内容)
- [🎯 快速开始](#-快速开始)
  - [HTTP服务器](#http服务器)
  - [JSON API](#json-api)
- [📖 系统文档](#-系统文档)
- [🔗 相关资源](#-相关资源)
- [🎓 学习路径](#-学习路径)
  - [初学者](#初学者)
  - [进阶者](#进阶者)
  - [专家](#专家)

## 📚 文档导航

### 核心内容

1. **[核心API参考](./01-核心API参考.md)** ⭐⭐⭐⭐⭐
   - net/http, encoding/json, context, sync
   - fmt, io, time, os, strings, strconv
   - errors, log
   - 10个核心包完整参考

2. **[标准库API](./02-标准库API.md)** ⭐⭐⭐⭐⭐
   - bufio, path/filepath, regexp
   - math, math/rand, sort, flag
   - database/sql, crypto, compress
   - 15+标准库包详解

3. **[常用第三方库](./03-常用第三方库.md)** ⭐⭐⭐⭐⭐
   - Web框架: Gin, Fiber, Echo
   - 数据库: GORM, sqlx
   - RPC: gRPC, Gorilla Mux
   - 缓存: Redis, etcd
   - 消息队列: Kafka, NATS
   - 测试: testify, gomock
   - 认证: JWT, bcrypt

4. **[API设计指南](./04-API设计指南.md)** ⭐⭐⭐⭐⭐
   - RESTful API设计原则
   - 错误处理与状态码
   - 版本控制策略
   - 安全性最佳实践
   - 分页、过滤、排序

---

## 🎯 快速开始

### HTTP服务器

```go
import "net/http"

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
})
http.ListenAndServe(":8080", nil)
```

### JSON API

```go
import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Message string `json:"message"`
}

http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Message: "Success"})
})
```

---

## 📖 系统文档

- **[知识图谱](./00-知识图谱.md)**: API知识体系全景图
- **[对比矩阵](./00-对比矩阵.md)**: 不同API方案对比
- **[概念定义体系](./00-概念定义体系.md)**: API核心概念详解

---

## 🔗 相关资源

- [Go标准库文档](https://pkg.go.dev/std)
- [Awesome Go](https://github.com/avelino/awesome-go)
- [Go Packages](https://pkg.go.dev/)

---

## 🎓 学习路径

### 初学者
1. [核心API参考](./01-核心API参考.md)
2. [标准库API](./02-标准库API.md)
3. 实践：构建简单HTTP服务

### 进阶者
1. [常用第三方库](./03-常用第三方库.md)
2. [API设计指南](./04-API设计指南.md)
3. 实践：构建RESTful API

### 专家
1. 深入标准库源码
2. 自定义第三方库
3. API性能优化

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3
