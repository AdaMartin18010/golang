# HTTP协议基础

> **简介**: 深入理解HTTP协议基础，包括报文结构、请求方法、状态码和HTTP/1.1与HTTP/2对比

> **版本**: Go 1.23+  
> **难度**: ⭐⭐⭐  
> **标签**: #Web #HTTP #协议 #网络编程

<!-- TOC START -->
- [HTTP协议基础](#http协议基础)
  - [📚 理论分析](#-理论分析)
    - [HTTP协议简介](#http协议简介)
    - [报文结构](#报文结构)
    - [常用HTTP方法](#常用http方法)
    - [状态码详解](#状态码详解)
    - [HTTP版本演进](#http版本演进)
  - [💻 Go语言HTTP实战](#-go语言http实战)
    - [HTTP客户端](#http客户端)
    - [HTTP服务器](#http服务器)
    - [中间件模式](#中间件模式)
  - [🔒 安全性考虑](#-安全性考虑)
  - [⚡ 性能优化](#-性能优化)
  - [🎯 最佳实践](#-最佳实践)
  - [🔍 常见问题](#-常见问题)
  - [📚 扩展阅读](#-扩展阅读)
<!-- TOC END -->


## 📋 目录


- [📚 理论分析](#-理论分析)
  - [HTTP协议简介](#http协议简介)
  - [报文结构](#报文结构)
    - [请求报文](#请求报文)
    - [响应报文](#响应报文)
  - [常用HTTP方法](#常用http方法)
  - [状态码详解](#状态码详解)
    - [1xx 信息响应](#1xx-信息响应)
    - [2xx 成功](#2xx-成功)
    - [3xx 重定向](#3xx-重定向)
    - [4xx 客户端错误](#4xx-客户端错误)
    - [5xx 服务器错误](#5xx-服务器错误)
  - [HTTP版本演进](#http版本演进)
    - [HTTP/1.0（1996）](#http101996)
    - [HTTP/1.1（1999）](#http111999)
    - [HTTP/2（2015）](#http22015)
    - [HTTP/3（2022）](#http32022)
- [💻 Go语言HTTP实战](#-go语言http实战)
  - [HTTP客户端](#http客户端)
  - [HTTP服务器](#http服务器)
  - [中间件模式](#中间件模式)
- [🔒 安全性考虑](#-安全性考虑)
  - [1. HTTPS](#1-https)
  - [2. CORS处理](#2-cors处理)
  - [3. 防止常见攻击](#3-防止常见攻击)
- [⚡ 性能优化](#-性能优化)
  - [1. 连接池配置](#1-连接池配置)
  - [2. HTTP/2支持](#2-http2支持)
  - [3. 响应缓存](#3-响应缓存)
- [🎯 最佳实践](#-最佳实践)
- [🔍 常见问题](#-常见问题)
- [📚 扩展阅读](#-扩展阅读)

## 📚 理论分析

### HTTP协议简介

**HTTP（HyperText Transfer Protocol）** 是Web通信的基础协议：

- **请求-响应模型**：客户端发送请求，服务器返回响应
- **无状态协议**：每个请求独立，服务器不保存客户端状态
- **基于文本**：HTTP/1.x使用文本格式（易读易调试）
- **应用层协议**：基于TCP/IP，默认端口80（HTTP）/443（HTTPS）

**核心特点：**
- 简单：请求-响应模式易于理解
- 可扩展：通过Headers添加元数据
- 无状态：通过Cookie/Session实现会话管理
- 灵活：支持多种内容类型（HTML、JSON、XML等）

### 报文结构

#### 请求报文

```
请求行（Request Line）
请求头（Headers）
空行
请求体（Body，可选）
```

**请求报文示例：**

```http
POST /api/users HTTP/1.1
Host: api.example.com
User-Agent: Go-http-client/1.1
Content-Type: application/json
Content-Length: 45
Authorization: Bearer eyJhbGc...
Accept: application/json

{"name":"John Doe","email":"john@example.com"}
```

#### 响应报文

```
状态行（Status Line）
响应头（Headers）
空行
响应体（Body）
```

**响应报文示例：**

```http
HTTP/1.1 201 Created
Content-Type: application/json; charset=UTF-8
Content-Length: 78
Date: Mon, 27 Oct 2025 12:00:00 GMT
Server: Go-Server/1.0

{"id":123,"name":"John Doe","email":"john@example.com","created_at":"2025-10-27T12:00:00Z"}
```

### 常用HTTP方法

| 方法 | 说明 | 幂等性 | 安全性 | 常见用途 |
|------|------|--------|--------|----------|
| **GET** | 获取资源 | ✅ | ✅ | 查询数据、获取页面 |
| **POST** | 提交数据 | ❌ | ❌ | 创建资源、提交表单 |
| **PUT** | 更新资源 | ✅ | ❌ | 完整更新资源 |
| **PATCH** | 部分更新 | ❌ | ❌ | 部分字段更新 |
| **DELETE** | 删除资源 | ✅ | ❌ | 删除数据 |
| **HEAD** | 获取元信息 | ✅ | ✅ | 检查资源是否存在 |
| **OPTIONS** | 查询支持的方法 | ✅ | ✅ | CORS预检请求 |
| **CONNECT** | 建立隧道 | ❌ | ❌ | HTTPS代理 |
| **TRACE** | 回显请求 | ✅ | ✅ | 调试（通常禁用） |

**幂等性说明：**
- **幂等**：多次执行产生相同结果（GET、PUT、DELETE）
- **非幂等**：多次执行结果不同（POST、PATCH）

### 状态码详解

#### 1xx 信息响应

- **100 Continue**: 客户端应继续发送请求体
- **101 Switching Protocols**: 切换协议（如WebSocket）

#### 2xx 成功

- **200 OK**: 请求成功
- **201 Created**: 资源创建成功
- **202 Accepted**: 请求已接受，但未完成
- **204 No Content**: 成功但无响应体
- **206 Partial Content**: 部分内容（断点续传）

#### 3xx 重定向

- **301 Moved Permanently**: 永久重定向
- **302 Found**: 临时重定向
- **304 Not Modified**: 资源未修改（缓存有效）
- **307 Temporary Redirect**: 临时重定向（保持方法）
- **308 Permanent Redirect**: 永久重定向（保持方法）

#### 4xx 客户端错误

- **400 Bad Request**: 请求格式错误
- **401 Unauthorized**: 未认证
- **403 Forbidden**: 无权限
- **404 Not Found**: 资源不存在
- **405 Method Not Allowed**: 方法不支持
- **409 Conflict**: 冲突（如并发修改）
- **429 Too Many Requests**: 请求过多（限流）

#### 5xx 服务器错误

- **500 Internal Server Error**: 服务器内部错误
- **502 Bad Gateway**: 网关错误
- **503 Service Unavailable**: 服务不可用
- **504 Gateway Timeout**: 网关超时

### HTTP版本演进

#### HTTP/1.0（1996）

- 每个请求一个TCP连接
- 无连接复用
- 无Host头（一个IP只能一个网站）

#### HTTP/1.1（1999）

```
优点：
✅ 持久连接（Keep-Alive）
✅ 管道化（Pipelining）
✅ 分块传输（Chunked Transfer）
✅ Host头支持虚拟主机
✅ 缓存控制增强

缺点：
❌ 队头阻塞（Head-of-Line Blocking）
❌ 串行请求（一个响应完成才能发下一个）
❌ Header冗余（每次重复发送）
```

#### HTTP/2（2015）

```
核心改进：
✅ 二进制分帧（Binary Framing）
✅ 多路复用（Multiplexing） - 一个连接并发多个请求
✅ Header压缩（HPACK）
✅ 服务器推送（Server Push）
✅ 流优先级（Stream Priority）

性能提升：
- 减少延迟50-70%
- 减少带宽10-30%
```

#### HTTP/3（2022）

```
基于QUIC协议（UDP）：
✅ 消除队头阻塞（TCP层）
✅ 连接迁移（IP变化不断连）
✅ 0-RTT连接建立
✅ 更好的拥塞控制
```

---

## 💻 Go语言HTTP实战

### HTTP客户端

**基本GET请求：**

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, err := http.Get("https://api.example.com/users")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()  // 必须关闭
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %d\n", resp.StatusCode)
    fmt.Printf("Body: %s\n", string(body))
}
```

**POST请求（JSON）：**

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(user User) error {
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    
    resp, err := http.Post(
        "https://api.example.com/users",
        "application/json",
        bytes.NewBuffer(data),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    return nil
}
```

**自定义HTTP Client（推荐）：**

```go
var client = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
    },
}

func makeRequest(url string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("User-Agent", "MyApp/1.0")
    req.Header.Set("Accept", "application/json")
    
    return client.Do(req)
}
```

### HTTP服务器

**基本服务器：**

```go
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // 设置响应头
    w.Header().Set("Content-Type", "application/json")
    
    // 写入响应
    fmt.Fprintf(w, `{"message":"Hello, World!","method":"%s"}`, r.Method)
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OK")
    })
    
    fmt.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}
```

**RESTful API示例：**

```go
type UserHandler struct {
    store *UserStore
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.getUsers(w, r)
    case http.MethodPost:
        h.createUser(w, r)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.store.List()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    
    if err := h.store.Create(&user); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

### 中间件模式

```go
// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// 认证中间件
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// 使用中间件
func main() {
    handler := http.HandlerFunc(myHandler)
    http.Handle("/api/", loggingMiddleware(authMiddleware(handler)))
    http.ListenAndServe(":8080", nil)
}
```

---

## 🔒 安全性考虑

### 1. HTTPS

```go
// 启动HTTPS服务器
func main() {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}
```

### 2. CORS处理

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### 3. 防止常见攻击

- **XSS**：输出编码，使用Content-Security-Policy
- **CSRF**：Token验证，SameSite Cookie
- **SQL注入**：参数化查询
- **DoS**：限流、超时设置

---

## ⚡ 性能优化

### 1. 连接池配置

```go
var client = &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,   // 总空闲连接数
        MaxIdleConnsPerHost: 10,    // 每个host空闲连接数
        MaxConnsPerHost:     50,    // 每个host最大连接数
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        DisableCompression:  false,  // 启用gzip
    },
    Timeout: 30 * time.Second,
}
```

### 2. HTTP/2支持

```go
import "golang.org/x/net/http2"

server := &http.Server{
    Addr:    ":8443",
    Handler: handler,
}

http2.ConfigureServer(server, &http2.Server{})
server.ListenAndServeTLS("cert.pem", "key.pem")
```

### 3. 响应缓存

```go
func cacheMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "public, max-age=3600")
        w.Header().Set("ETag", generateETag(r.URL.Path))
        next.ServeHTTP(w, r)
    })
}
```

---

## 🎯 最佳实践

1. **始终关闭Response.Body**
   ```go
   resp, err := http.Get(url)
   if err != nil {
       return err
   }
   defer resp.Body.Close()  // 必须！
   ```

2. **设置合理的超时**
   ```go
   client.Timeout = 30 * time.Second
   ```

3. **复用HTTP Client**
   - 不要每次请求创建新Client
   - Client是并发安全的

4. **正确处理状态码**
   ```go
   if resp.StatusCode != http.StatusOK {
       return fmt.Errorf("unexpected status: %d", resp.StatusCode)
   }
   ```

5. **使用Context控制超时**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   
   req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
   resp, err := client.Do(req)
   ```

---

## 🔍 常见问题

**Q: HTTP是有状态的吗？**
A: HTTP本身无状态，但可通过Cookie/Session/JWT等机制实现会话管理。

**Q: GET和POST的区别？**
A: 
- GET：查询数据，参数在URL，幂等，可缓存
- POST：提交数据，参数在Body，非幂等，不可缓存

**Q: 何时使用PUT vs PATCH？**
A:
- PUT：完整替换资源（幂等）
- PATCH：部分更新资源（非幂等）

**Q: HTTP/2一定比HTTP/1.1快吗？**
A: 通常是，但在高延迟或小文件场景下优势不明显。

**Q: 如何实现文件上传？**
A:
```go
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // 处理file
}
```

---

## 📚 扩展阅读

- [MDN HTTP协议详解](https://developer.mozilla.org/zh-CN/docs/Web/HTTP)
- [RFC 7230-7235: HTTP/1.1规范](https://tools.ietf.org/html/rfc7230)
- [RFC 7540: HTTP/2规范](https://tools.ietf.org/html/rfc7540)
- [Go net/http官方文档](https://golang.org/pkg/net/http/)
- [《HTTP权威指南》](https://www.oreilly.com/library/view/http-the-definitive/1565925092/)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月27日  
**文档状态**: 已优化  
**适用版本**: Go 1.25.3+
