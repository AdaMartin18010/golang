# TS-NET-004: WebSocket in Go

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #websocket #realtime #golang #gorilla #messaging
> **权威来源**:
>
> - [WebSocket RFC 6455](https://tools.ietf.org/html/rfc6455) - IETF
> - [Gorilla WebSocket](https://github.com/gorilla/websocket) - Gorilla toolkit

---

## 1. WebSocket Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       WebSocket Architecture                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Connection Upgrade:                                                         │
│  ┌─────────────┐                    ┌─────────────┐                         │
│  │   Client    │────────────────────│   Server    │                         │
│  └──────┬──────┘  HTTP Upgrade      └──────┬──────┘                         │
│         │                                  │                                  │
│  GET /ws HTTP/1.1                         HTTP/1.1 101 Switching Protocols   │
│  Host: server.com                         Upgrade: websocket                 │
│  Upgrade: websocket                       Connection: Upgrade                │
│  Connection: Upgrade                      Sec-WebSocket-Accept: s3pPLMBiTx   │
│  Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==                                 │
│  Sec-WebSocket-Version: 13                                                   │
│                                                                              │
│  After Upgrade:                                                              │
│  ┌─────────────┐      Bidirectional         ┌─────────────┐                 │
│  │   Client    │◄───────────────────────────►│   Server    │                 │
│  │   (Browser) │      Full-Duplex            │    (Go)     │                 │
│  └─────────────┘                             └─────────────┘                 │
│                                                                              │
│  Frame Structure:                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  0                   1                   2                   3       │   │
│  │  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1   │   │
│  │ +-+-+-+-+-------+-+-------------+-------------------------------+   │   │
│  │ |F|R|R|R| opcode|M| Payload len |    Extended payload length    |   │   │
│  │ |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |   │   │
│  │ |N|V|V|V|       |S|             |   (if payload len==126/127)   |   │   │
│  │ | |1|2|3|       |K|             |                               |   │   │
│  │ +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +   │   │
│  │ |     Extended payload length continued, if payload len == 127  |   │   │
│  │ + - - - - - - - - - - - - - - - +-------------------------------+   │   │
│  │ |                               |Masking-key, if MASK set to 1  |   │   │
│  │ +-------------------------------+-------------------------------+   │   │
│  │ | Masking-key (continued)       |          Payload Data         |   │   │
│  │ +-------------------------------- - - - - - - - - - - - - - - - +   │   │
│  │ :                     Payload Data continued ...                :   │   │
│  │ + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +   │   │
│  │ |                     Payload Data continued ...                |   │   │
│  │ +---------------------------------------------------------------+   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. WebSocket Server with Gorilla

```go
package main

import (
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow all origins in development
        // In production, check against allowed origins
        return true
    },
}

// Hub maintains the set of active clients
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

func newHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(512 * 1024) // 512KB max message size
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        // Process message
        c.hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued messages
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}

func main() {
    hub := newHub()
    go hub.run()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 3. WebSocket Client

```go
package main

import (
    "log"
    "net/url"
    "os"
    "os/signal"
    "time"

    "github.com/gorilla/websocket"
)

func client() {
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
    log.Printf("connecting to %s", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("dial:", err)
    }
    defer c.Close()

    done := make(chan struct{})

    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("read:", err)
                return
            }
            log.Printf("recv: %s", message)
        }
    }()

    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-done:
            return
        case t := <-ticker.C:
            err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
            if err != nil {
                log.Println("write:", err)
                return
            }
        case <-interrupt:
            log.Println("interrupt")

            // Clean close
            err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            if err != nil {
                log.Println("write close:", err)
                return
            }
            select {
            case <-done:
            case <-time.After(time.Second):
            }
            return
        }
    }
}
```

---

## 4. Checklist

```
WebSocket Checklist:
□ Use Gorilla WebSocket library
□ Handle connection upgrades properly
□ Implement ping/pong for keepalive
□ Handle concurrent access with mutex
□ Set read/write deadlines
□ Limit message size
□ Handle graceful shutdown
□ Implement reconnection logic
□ Monitor connection health
```
