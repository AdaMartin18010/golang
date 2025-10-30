# WebSocket 实时通信

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.23+

---

> **难度**: ⭐⭐⭐⭐
> **标签**: #Web #WebSocket #实时通信 #gorilla

## 📋 目录

- [📚 **理论分析**](#理论分析)
  - [**WebSocket原理**](#websocket原理)
  - [**协议流程**](#协议流程)
  - [**安全与性能**](#安全与性能)
- [💻 **代码示例**](#代码示例)
  - [**标准库+第三方库（gorilla/websocket）**](#标准库+第三方库gorillawebsocket)
  - [**Gin集成WebSocket**](#gin集成websocket)
  - [**Echo集成WebSocket**](#echo集成websocket)
  - [**Fiber集成WebSocket**](#fiber集成websocket)
- [🧪 **测试代码**](#测试代码)
- [🎯 **最佳实践**](#最佳实践)
- [🔍 **常见问题**](#常见问题)
- [📚 **扩展阅读**](#扩展阅读)

## 📚 **理论分析**

### **WebSocket原理**

- WebSocket是一种全双工、持久化的网络通信协议，基于TCP。
- 通过HTTP/1.1升级握手（Upgrade: websocket），建立后可双向实时通信。
- 适合聊天室、实时推送、在线协作等场景。

### **协议流程**

- 客户端发起HTTP请求，包含`Upgrade: websocket`头
- 服务器响应101 Switching Protocols，升级为WebSocket
- 后续数据以帧（frame）方式双向传输

### **安全与性能**

- 建议使用wss（TLS加密）
- 需做连接数、消息大小、心跳检测等限制

## 💻 **代码示例**

### **标准库+第三方库（gorilla/websocket）**

```go
package main
import (
    "github.com/gorilla/websocket"
    "net/http"
)
var upgrader = websocket.Upgrader{}
func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil { return }
    defer conn.Close()
    for {
        mt, msg, err := conn.ReadMessage()
        if err != nil { break }
        conn.WriteMessage(mt, msg) // echo
    }
}
func main() {
    http.HandleFunc("/ws", wsHandler)
    http.ListenAndServe(":8080", nil)
}
```

### **Gin集成WebSocket**

```go
package main
import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "net/http"
)
var upgrader = websocket.Upgrader{}
func main() {
    r := gin.Default()
    r.GET("/ws", func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil { return }
        defer conn.Close()
        for {
            mt, msg, err := conn.ReadMessage()
            if err != nil { break }
            conn.WriteMessage(mt, msg)
        }
    })
    r.Run(":8080")
}
```

### **Echo集成WebSocket**

```go
package main
import (
    "github.com/labstack/echo/v4"
    "github.com/gorilla/websocket"
    "net/http"
)
var upgrader = websocket.Upgrader{}
func main() {
    e := echo.New()
    e.GET("/ws", func(c echo.Context) error {
        conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
        if err != nil { return err }
        defer conn.Close()
        for {
            mt, msg, err := conn.ReadMessage()
            if err != nil { break }
            conn.WriteMessage(mt, msg)
        }
        return nil
    })
    e.Logger.Fatal(e.Start(":8080"))
}
```

### **Fiber集成WebSocket**

```go
package main
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
)
func main() {
    app := fiber.New()
    app.Use("/ws", websocket.New(func(c *websocket.Conn) {
        defer c.Close()
        for {
            mt, msg, err := c.ReadMessage()
            if err != nil { break }
            c.WriteMessage(mt, msg)
        }
    }))
    app.Listen(":8080")
}
```

## 🧪 **测试代码**

```go
// 可用websocket客户端或浏览器测试
```

## 🎯 **最佳实践**

- 限制最大连接数和消息大小，防止滥用
- 定期心跳检测，及时断开无效连接
- 生产环境建议用nginx反向代理wss
- 合理处理异常和断线重连

## 🔍 **常见问题**

- Q: WebSocket和HTTP的区别？
  A: WebSocket为持久化双向通信，HTTP为短连接请求-响应
- Q: 如何做身份认证？
  A: 握手时校验token或cookie
- Q: 如何广播消息？
  A: 维护连接池，遍历发送

## 📚 **扩展阅读**

- [gorilla/websocket文档](https://pkg.go.dev/github.com/gorilla/websocket)
- [MDN WebSocket协议](https://developer.mozilla.org/zh-CN/docs/Web/API/WebSockets_API)
- [Fiber WebSocket文档](https://docs.gofiber.io/api/websocket)

---

**文档维护者**: Go Documentation Team
**最后更新**: 2025-10-29
**文档状态**: 完成
**适用版本**: Go 1.25.3+
