# URL工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [URL工具](#url工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
    - [2.1 URL解析和构建](#21-url解析和构建)
    - [2.2 查询参数操作](#22-查询参数操作)
    - [2.3 URL组件操作](#23-url组件操作)
    - [2.4 URL验证和转换](#24-url验证和转换)
    - [2.5 URL编码](#25-url编码)
    - [2.6 URL安全](#26-url安全)
  - [3. 使用示例](#3-使用示例)
    - [3.1 URL构建](#31-url构建)
    - [3.2 查询参数操作](#32-查询参数操作)
    - [3.3 URL组件操作](#33-url组件操作)
    - [3.4 URL验证](#34-url验证)
    - [3.5 URL编码](#35-url编码)
    - [3.6 URL安全](#36-url安全)

---

## 1. 概述

URL工具提供了丰富的URL操作函数，简化常见的URL处理任务。

---

## 2. 功能特性

### 2.1 URL解析和构建

- `Parse`: 解析URL
- `ParseRequestURI`: 解析请求URI
- `BuildURL`: 构建URL
- `JoinPath`: 连接路径
- `Resolve`: 解析相对URL

### 2.2 查询参数操作

- `AddQuery`: 添加查询参数
- `AddQueries`: 批量添加查询参数
- `RemoveQuery`: 移除查询参数
- `GetQuery`: 获取查询参数值
- `GetAllQueries`: 获取所有查询参数
- `BuildQueryString`: 构建查询字符串
- `ParseQueryString`: 解析查询字符串

### 2.3 URL组件操作

- `SetScheme`: 设置URL协议
- `SetHost`: 设置URL主机
- `SetPath`: 设置URL路径
- `GetScheme`: 获取URL协议
- `GetHost`: 获取URL主机
- `GetPath`: 获取URL路径
- `GetDomain`: 获取域名
- `GetPort`: 获取端口号

### 2.4 URL验证和转换

- `IsValid`: 检查URL是否有效
- `IsAbsolute`: 检查URL是否为绝对路径
- `Normalize`: 规范化URL
- `IsHTTPS`: 检查是否为HTTPS
- `IsHTTP`: 检查是否为HTTP
- `ToHTTPS`: 转换为HTTPS
- `ToHTTP`: 转换为HTTP

### 2.5 URL编码

- `Encode`: 编码URL
- `Decode`: 解码URL

### 2.6 URL安全

- `MaskURL`: 掩码URL（隐藏敏感信息）

---

## 3. 使用示例

### 3.1 URL构建

```go
import "github.com/yourusername/golang/pkg/utils/url"

// 构建URL
result, err := url.BuildURL("https://api.example.com", "/users", map[string]string{
    "page":  "1",
    "limit": "10",
})
// 结果: https://api.example.com/users?limit=10&page=1
```

### 3.2 查询参数操作

```go
// 添加查询参数
result, err := url.AddQuery("https://api.example.com/users", "page", "1")

// 批量添加查询参数
params := map[string]string{"page": "1", "limit": "10"}
result, err := url.AddQueries("https://api.example.com/users", params)

// 获取查询参数
value, err := url.GetQuery("https://api.example.com/users?page=1", "page")

// 移除查询参数
result, err := url.RemoveQuery("https://api.example.com/users?page=1", "page")
```

### 3.3 URL组件操作

```go
// 获取域名
domain, err := url.GetDomain("https://api.example.com:8080/users")

// 获取端口
port, err := url.GetPort("https://api.example.com:8080/users")

// 设置协议
result, err := url.SetScheme("http://api.example.com", "https")
```

### 3.4 URL验证

```go
// 检查URL是否有效
if url.IsValid("https://api.example.com") {
    // URL有效
}

// 检查是否为绝对路径
if url.IsAbsolute("https://api.example.com") {
    // 绝对路径
}

// 检查是否为HTTPS
if url.IsHTTPS("https://api.example.com") {
    // HTTPS URL
}
```

### 3.5 URL编码

```go
// 编码URL
encoded := url.Encode("hello world") // "hello%20world"

// 解码URL
decoded, err := url.Decode("hello%20world") // "hello world"
```

### 3.6 URL安全

```go
// 掩码URL（隐藏敏感信息）
masked, err := url.MaskURL("https://user:pass@api.example.com/users?token=secret123")
// 结果: https://api.example.com/users?token=***
```

---

**更新日期**: 2025-11-11
