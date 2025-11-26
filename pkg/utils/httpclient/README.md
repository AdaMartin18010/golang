# HTTP客户端工具

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [HTTP客户端工具](#http客户端工具)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
  - [2. 功能特性](#2-功能特性)
  - [3. 使用示例](#3-使用示例)
  - [4. 最佳实践](#4-最佳实践)

---

## 1. 概述

HTTP客户端工具提供了简单易用的HTTP客户端，简化HTTP请求的发送和处理。

---

## 2. 功能特性

### 2.1 基础功能

- GET、POST、PUT、DELETE、PATCH请求支持
- 请求参数和查询参数支持
- 请求头管理
- 超时控制
- Context支持

### 2.2 响应处理

- JSON响应解析
- 响应状态码检查
- 响应头访问
- 响应体字符串转换

---

## 3. 使用示例

### 3.1 创建客户端

```go
import "github.com/yourusername/golang/pkg/utils/httpclient"

client := httpclient.NewClient(httpclient.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
})
```

### 3.2 GET请求

```go
// 带查询参数
params := map[string]string{
    "page": "1",
    "limit": "10",
}
resp, err := client.Get(ctx, "/users", params)
if err != nil {
    // 处理错误
}

// 解析JSON响应
var users []User
err = resp.JSON(&users)
```

### 3.3 POST请求

```go
body := map[string]interface{}{
    "name":  "John",
    "email": "john@example.com",
}

resp, err := client.Post(ctx, "/users", body, nil)
if err != nil {
    // 处理错误
}

if resp.IsSuccess() {
    var user User
    resp.JSON(&user)
}
```

### 3.4 PUT请求

```go
body := map[string]interface{}{
    "name": "John Updated",
}

resp, err := client.Put(ctx, "/users/1", body, nil)
```

### 3.5 DELETE请求

```go
resp, err := client.Delete(ctx, "/users/1", nil)
```

### 3.6 设置请求头

```go
// 设置默认请求头
client.SetHeader("Authorization", "Bearer token")

// 单次请求设置请求头
headers := map[string]string{
    "X-Custom-Header": "value",
}
resp, err := client.Post(ctx, "/users", body, headers)
```

### 3.7 使用默认客户端

```go
// 使用默认客户端发送请求
resp, err := httpclient.Get(ctx, "https://api.example.com/users", nil)
resp, err := httpclient.Post(ctx, "https://api.example.com/users", body)
```

### 3.8 响应处理

```go
resp, err := client.Get(ctx, "/users", nil)
if err != nil {
    // 处理错误
}

// 检查状态码
if resp.IsSuccess() {
    // 处理成功响应
}

if resp.IsError() {
    // 处理错误响应
}

// 获取响应头
contentType := resp.GetHeader("Content-Type")

// 获取响应体字符串
bodyStr := resp.String()

// 解析JSON
var data map[string]interface{}
err = resp.JSON(&data)
```

---

## 4. 最佳实践

### 4.1 超时设置

设置合理的超时时间，避免请求长时间阻塞：

```go
client := httpclient.NewClient(httpclient.Config{
    Timeout: 10 * time.Second, // 根据实际情况设置
})
```

### 4.2 Context使用

使用Context控制请求取消和超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Get(ctx, "/users", nil)
```

### 4.3 错误处理

始终检查错误和响应状态码：

```go
resp, err := client.Get(ctx, "/users", nil)
if err != nil {
    // 处理网络错误
    return err
}

if resp.IsError() {
    // 处理HTTP错误
    return fmt.Errorf("HTTP error: %d", resp.StatusCode)
}
```

### 4.4 请求头管理

使用SetHeader设置通用请求头，使用headers参数设置特定请求头：

```go
// 设置通用请求头
client.SetHeader("Authorization", "Bearer token")

// 特定请求使用不同请求头
headers := map[string]string{
    "X-Request-ID": "unique-id",
}
resp, err := client.Post(ctx, "/users", body, headers)
```

---

**更新日期**: 2025-11-11
