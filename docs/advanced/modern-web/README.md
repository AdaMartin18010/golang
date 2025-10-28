# Go现代Web开发

Go现代Web开发完整指南，涵盖现代Web框架、实时通信、GraphQL和微服务网关。

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

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**最后更新**: 2025-10-28
