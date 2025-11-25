# Go HTTP服务器进阶

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.23+

---

## 📋 目录

- [Go HTTP服务器进阶](#go-http服务器进阶)
  - [📋 目录](#-目录)
  - [📚 **理论分析**](#-理论分析)
    - [**HTTP服务器核心原理**](#http服务器核心原理)
    - [**服务器配置参数**](#服务器配置参数)
    - [**优雅关闭与重启**](#优雅关闭与重启)
    - [**并发与性能优化**](#并发与性能优化)
  - [💻 **代码示例**](#-代码示例)
    - [**自定义HTTP服务器与超时配置**](#自定义http服务器与超时配置)
    - [**优雅关闭服务器**](#优雅关闭服务器)
    - [**静态文件服务**](#静态文件服务)
  - [🧪 **测试代码**](#-测试代码)
  - [🎯 **最佳实践**](#-最佳实践)
  - [🔍 **常见问题**](#-常见问题)
  - [📚 **扩展阅读**](#-扩展阅读)

---

## 📚 **理论分析**

### **HTTP服务器核心原理**

- Go内置`http.Server`类型，支持高并发、可配置、易扩展。
- 支持HTTP/1.1和HTTP/2，自动处理连接复用。
- 服务器可通过自定义`Handler`、`ServeMux`实现复杂路由。

### **服务器配置参数**

- `Addr`：监听地址（如":8080"）
- `Handler`：请求处理器
- `ReadTimeout`/`WriteTimeout`：超时控制，防止慢连接攻击
- `IdleTimeout`：空闲连接超时
- `TLSConfig`：HTTPS支持

### **优雅关闭与重启**

- 使用`http.Server.Shutdown(ctx)`实现平滑关闭，确保所有连接处理完毕
- 可结合`signal.Notify`监听系统信号（如SIGINT、SIGTERM）

### **并发与性能优化**

- 每个请求由独立Goroutine处理，天然高并发
- 可通过连接池、对象池、限流等手段提升性能
- 静态文件服务建议用`http.ServeContent`或`http.FileServer`

## 💻 **代码示例**

### **自定义HTTP服务器与超时配置**

```go
package main
import (
    "fmt"
    "net/http"
    "time"
)
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, Custom Server!")
}
func main() {
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      http.HandlerFunc(hello),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    fmt.Println("Server running on :8080")
    if err := srv.ListenAndServe(); err != nil {
        fmt.Println("Server stopped:", err)
    }
}
```

### **优雅关闭服务器**

```go
package main
import (
    "Context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)
func main() {
    srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
    go func() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "Graceful shutdown demo")
        })
        srv.ListenAndServe()
    }()
    quit := make(Channel os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    ctx, cancel := Context.WithTimeout(Context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        fmt.Println("Shutdown error:", err)
    }
    fmt.Println("Server gracefully stopped")
}
```

### **静态文件服务**

```go
package main
import "net/http"
func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.ListenAndServe(":8080", nil)
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
func TestCustomServer(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Test OK"))
    }).ServeHTTP(w, req)
    if w.Body.String() != "Test OK" {
        t.Errorf("unexpected response: %s", w.Body.String())
    }
}
```

## 🎯 **最佳实践**

- 合理设置超时，防止慢连接攻击
- 使用优雅关闭，避免强制中断请求
- 静态文件服务建议隔离目录，防止目录遍历
- 日志与监控集成，便于排查问题
- 生产环境建议使用反向代理（如Nginx）

## 🔍 **常见问题**

- Q: 如何支持HTTPS？
  A: 使用`srv.ListenAndServeTLS(cert, key)`
- Q: 如何限制最大并发连接数？
  A: 可用第三方库或自定义连接池
- Q: 如何实现健康检查？
  A: 提供`/healthz`等接口，返回200状态码

## 📚 **扩展阅读**
