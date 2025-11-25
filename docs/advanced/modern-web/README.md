# Go现代Web开发

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go现代Web开发](#go现代web开发)
  - [� 目录](#-目录)
  - [📚 核心内容](#-核心内容)
  - [🚀 WebSocket示例](#-websocket示例)
  - [📖 系统文档](#-系统文档)

---

## 📚 核心内容

1. **[现代Web框架](./01-现代Web框架.md)** ⭐⭐⭐⭐
2. **[实时通信](./02-实时通信.md)** ⭐⭐⭐⭐⭐
   - WebSocket
   - Server-Sent Events (SSE)
3. **[GraphQL](./03-GraphQL.md)** ⭐⭐⭐⭐⭐
4. **[微服务网关](./04-微服务网关.md)** ⭐⭐⭐⭐⭐
5. **[服务网格](./05-服务网格.md)** ⭐⭐⭐⭐
6. **[云原生实践](./06-云原生实践.md)** ⭐⭐⭐⭐⭐

---

## 🚀 WebSocket示例

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            return
        }
        conn.WriteMessage(messageType, p)
    }
}
```

---

## 📖 系统文档
