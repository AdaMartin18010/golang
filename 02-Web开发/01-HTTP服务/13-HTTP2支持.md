# HTTP/2 支持

## 📚 **理论分析**

### **HTTP/2协议原理**

- HTTP/2是HTTP/1.1的升级版，采用二进制分帧协议，支持多路复用、头部压缩、服务器推送等特性。
- 解决了HTTP/1.1队头阻塞、连接复用不足等问题。
- 默认使用TLS（h2协议），更安全高效。

### **核心特性**

- 多路复用：单连接并发多请求，提升吞吐量
- 头部压缩：减少带宽消耗
- 服务器推送：主动推送资源
- 流量控制与优先级

## 💻 **Go语言HTTP/2实现**

- Go 1.6+标准库`net/http`自动支持HTTP/2（HTTPS下）
- 明确监听h2协议或用`golang.org/x/net/http2`包自定义

### **标准库自动支持**

```go
package main
import (
    "fmt"
    "net/http"
)
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Proto: %s", r.Proto)
    })
    // 需提供证书
    http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
}
```

### **自定义HTTP/2服务器**

```go
package main
import (
    "fmt"
    "net/http"
    "golang.org/x/net/http2"
)
func main() {
    srv := &http.Server{Addr: ":8443", Handler: http.DefaultServeMux}
    http2.ConfigureServer(srv, &http2.Server{})
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Proto: %s", r.Proto)
    })
    srv.ListenAndServeTLS("cert.pem", "key.pem")
}
```

### **Gin/Echo/Fiber等框架**

- 只需用`ListenAndServeTLS`启动，框架自动支持HTTP/2

## 📊 **性能与安全**

- 多路复用显著提升高并发场景性能
- 建议开启TLS，防止中间人攻击
- 合理配置流量控制，防止资源滥用

## 🎯 **最佳实践**

- 使用HTTPS，自动启用HTTP/2
- 配置合适的证书和加密套件
- 监控协议协商和性能指标
- 兼容HTTP/1.1客户端

## 🔍 **常见问题**

- Q: HTTP/2和HTTP/1.1兼容吗？
  A: 完全兼容，自动协商
- Q: 如何判断客户端是否用HTTP/2？
  A: 检查`r.Proto`字段
- Q: 服务器推送如何用？
  A: 标准库暂不支持，需第三方库

## 📚 **扩展阅读**

- [Go官方文档-HTTP/2](https://golang.org/pkg/net/http/#hdr-HTTP_2_Support)
- [HTTP/2 RFC7540](https://datatracker.ietf.org/doc/html/rfc7540)
- [MDN HTTP/2](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Overview#http2)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成
