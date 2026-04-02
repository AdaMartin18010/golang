# WebSocket

> **分类**: 开源技术堆栈

---

## gorilla/websocket

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
```

---

## 服务端

```go
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            break
        }

        // Echo
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            break
        }
    }
}
```

---

## 客户端

```go
c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
if err != nil {
    log.Fatal(err)
}
defer c.Close()

// 发送
c.WriteMessage(websocket.TextMessage, []byte("hello"))

// 接收
_, message, _ := c.ReadMessage()
fmt.Println(string(message))
```

---

## Hub 模式

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}
```
