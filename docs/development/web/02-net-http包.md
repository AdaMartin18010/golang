# Go标准库 net/http 包

> **简介**: 深入解析Go标准库net/http包，掌握HTTP服务器和客户端开发的核心技术
> **版本**: Go 1.23+  
> **难度**: ⭐⭐⭐  
> **标签**: #Web #net/http #HTTP服务器 #标准库

<!-- TOC START -->
- [Go标准库 net/http 包](#go标准库-nethttp-包)
  - [📚 **理论分析**](#-理论分析)
    - [**net/http包简介**](#nethttp包简介)
    - [**核心类型与接口**](#核心类型与接口)
    - [**路由与中间件机制**](#路由与中间件机制)
  - [💻 **代码示例**](#-代码示例)
    - [**最小HTTP服务器**](#最小http服务器)
    - [**自定义路由与多路复用**](#自定义路由与多路复用)
    - [**中间件实现（日志）**](#中间件实现日志)
    - [**HTTP客户端请求**](#http客户端请求)
  - [🧪 **测试代码**](#-测试代码)
  - [🎯 **最佳实践**](#-最佳实践)
  - [🔍 **常见问题**](#-常见问题)
  - [📚 **扩展阅读**](#-扩展阅读)
<!-- TOC END -->


## 📋 目录


- [📚 **理论分析**](#-理论分析)
  - [**net/http包简介**](#nethttp包简介)
  - [**核心类型与接口**](#核心类型与接口)
  - [**路由与中间件机制**](#路由与中间件机制)
- [💻 **代码示例**](#-代码示例)
  - [**最小HTTP服务器**](#最小http服务器)
  - [**自定义路由与多路复用**](#自定义路由与多路复用)
  - [**中间件实现（日志）**](#中间件实现日志)
  - [**HTTP客户端请求**](#http客户端请求)
- [🧪 **测试代码**](#-测试代码)
- [🎯 **最佳实践**](#-最佳实践)
- [🔍 **常见问题**](#-常见问题)
- [📚 **扩展阅读**](#-扩展阅读)

## 📚 **理论分析**

### **net/http包简介**

- Go内置的Web开发标准库，支持HTTP/1.1和HTTP/2。
- 提供高效、易用的服务端与客户端API。
- 支持路由、中间件、文件服务、表单处理等常见Web开发需求。

### **核心类型与接口**

- `http.Server`：HTTP服务器
- `http.Request`：请求对象
- `http.ResponseWriter`：响应写入接口
- `http.Handler`：处理器接口（`ServeHTTP(w, r)`）
- `http.Client`：HTTP客户端

### **路由与中间件机制**

- 路由通过`http.HandleFunc`或自定义`ServeMux`实现
- 中间件可通过包装`Handler`实现链式调用

## 💻 **代码示例**

### **最小HTTP服务器**

```go
package main
import (
    "fmt"
    "net/http"
)
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })
    http.ListenAndServe(":8080", nil)
}
```

### **自定义路由与多路复用**

```go
package main
import (
    "fmt"
    "net/http"
)
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello!")
}
func about(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "About page")
}
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    mux.HandleFunc("/about", about)
    http.ListenAndServe(":8080", mux)
}
```

### **中间件实现（日志）**

```go
package main
import (
    "log"
    "net/http"
    "time"
)
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}
func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
}
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    logged := logging(mux)
    http.ListenAndServe(":8080", logged)
}
```

### **HTTP客户端请求**

```go
package main
import (
    "fmt"
    "io/ioutil"
    "net/http"
)
func main() {
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

## 🧪 **测试代码**

```go
package main
import (
    "net/http"
    "net/http/httptest"
    "testing"
)
func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, Test!"))
    }).ServeHTTP(w, req)
    if w.Body.String() != "Hello, Test!" {
        t.Errorf("unexpected response: %s", w.Body.String())
    }
}
```

## 🎯 **最佳实践**

- 合理设置超时（`Server.ReadTimeout`/`WriteTimeout`）
- 使用`Context`管理请求生命周期
- 日志与错误处理分离
- 路由建议用第三方库（如Gin）做复杂项目
- 静态文件服务用`http.FileServer`

## 🔍 **常见问题**

- Q: 如何优雅关闭服务器？
  A: 使用`http.Server.Shutdown(ctx)`
- Q: 如何处理大文件上传？
  A: 设置`MaxBytesReader`限制，分块处理
- Q: 如何实现RESTful API？
  A: 结合路由、中间件、JSON序列化

## 📚 **扩展阅读**

- [Go官方文档-net/http](https://golang.org/pkg/net/http/)
- [Go by Example: HTTP Servers](https://gobyexample.com/http-servers)
- [Go by Example: HTTP Clients](https://gobyexample.com/http-clients)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025年10月20日  
**文档状态**: 完成  
**适用版本**: Go 1.25.3+
