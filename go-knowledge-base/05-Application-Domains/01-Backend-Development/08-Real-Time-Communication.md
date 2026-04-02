# 实时通信 (Real-Time Communication)

> **分类**: 成熟应用领域

---

## WebSocket Hub

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        c.hub.broadcast <- message
    }
}
```

---

## SSE (Server-Sent Events)

```go
func SSEHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }
    
    for {
        select {
        case <-r.Context().Done():
            return
        case msg := <-messages:
            fmt.Fprintf(w, "data: %s\n\n", msg)
            flusher.Flush()
        }
    }
}
```

---

## 长轮询

```go
func LongPollingHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
    defer cancel()
    
    select {
    case data := <-dataChannel:
        json.NewEncoder(w).Encode(data)
    case <-ctx.Done():
        w.WriteHeader(http.StatusNoContent)
    }
}
```

---

## 技术选择

| 技术 | 方向 | 适用场景 |
|------|------|----------|
| WebSocket | 双向 | 聊天、游戏 |
| SSE | 单向(服务器→客户端) | 股票、通知 |
| 长轮询 | 单向 | 简单实时 |
