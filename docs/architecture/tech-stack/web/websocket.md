# 1. 🔌 WebSocket 深度解析

> **简介**: 本文档详细阐述了 WebSocket 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 🔌 WebSocket 深度解析](#1--websocket-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 WebSocket 服务器](#131-websocket-服务器)
    - [1.3.2 连接管理](#132-连接管理)
    - [1.3.3 消息处理](#133-消息处理)
    - [1.3.4 心跳检测](#134-心跳检测)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 WebSocket 设计最佳实践](#141-websocket-设计最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**WebSocket 是什么？**

WebSocket 是一种在单个 TCP 连接上进行全双工通信的协议，支持实时双向数据传输。

**核心特性**:

- ✅ **全双工通信**: 客户端和服务器可以同时发送数据
- ✅ **低延迟**: 比 HTTP 轮询延迟更低
- ✅ **实时性**: 支持实时数据推送
- ✅ **标准协议**: 符合 RFC 6455 标准

---

## 1.2 选型论证

**为什么选择 WebSocket？**

**论证矩阵**:

| 评估维度 | 权重 | WebSocket | HTTP 轮询 | Server-Sent Events | gRPC Streaming | 说明 |
|---------|------|-----------|-----------|-------------------|----------------|------|
| **实时性** | 30% | 10 | 5 | 8 | 9 | WebSocket 实时性最好 |
| **双向通信** | 25% | 10 | 5 | 3 | 10 | WebSocket 支持双向 |
| **性能** | 20% | 9 | 5 | 8 | 9 | WebSocket 性能优秀 |
| **易用性** | 15% | 8 | 10 | 9 | 7 | WebSocket 易用性好 |
| **浏览器支持** | 10% | 10 | 10 | 9 | 6 | WebSocket 浏览器支持好 |
| **加权总分** | - | **9.30** | 6.00 | 7.50 | 8.40 | WebSocket 得分最高 |

**核心优势**:

1. **实时性（权重 30%）**:
   - 低延迟，实时推送
   - 适合实时应用场景
   - 比 HTTP 轮询效率高

2. **双向通信（权重 25%）**:
   - 客户端和服务器都可以主动发送
   - 适合交互式应用
   - 支持复杂通信模式

---

## 1.3 实际应用

### 1.3.1 WebSocket 服务器

**创建 WebSocket 服务器**:

```go
// internal/infrastructure/websocket/server.go
package websocket

import (
    "github.com/gorilla/websocket"
    "net/http"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 生产环境需要验证 Origin
    },
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error("Failed to upgrade connection", "error", err)
        return
    }
    defer conn.Close()

    // 处理连接
    handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
    for {
        // 读取消息
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            logger.Error("Failed to read message", "error", err)
            break
        }

        // 处理消息
        response := processMessage(message)

        // 发送响应
        if err := conn.WriteMessage(messageType, response); err != nil {
            logger.Error("Failed to write message", "error", err)
            break
        }
    }
}
```

### 1.3.2 连接管理

**连接 Hub 管理**:

```go
// Hub 管理所有 WebSocket 连接
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

type Client struct {
    hub    *Hub
    conn   *websocket.Conn
    send   chan []byte
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
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

// 广播消息
func (h *Hub) Broadcast(message []byte) {
    h.broadcast <- message
}
```

### 1.3.3 消息处理

**消息处理**:

```go
// 消息类型定义
type Message struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// 处理不同类型的消息
func processMessage(data []byte) []byte {
    var msg Message
    if err := json.Unmarshal(data, &msg); err != nil {
        return createErrorResponse("Invalid message format")
    }

    switch msg.Type {
    case "ping":
        return createResponse("pong", nil)
    case "subscribe":
        return handleSubscribe(msg.Payload)
    case "unsubscribe":
        return handleUnsubscribe(msg.Payload)
    default:
        return createErrorResponse("Unknown message type")
    }
}

func createResponse(msgType string, payload interface{}) []byte {
    msg := Message{
        Type:    msgType,
        Payload: payload,
    }
    data, _ := json.Marshal(msg)
    return data
}
```

### 1.3.4 心跳检测

**心跳检测**:

```go
// 心跳检测
const (
    pingPeriod = 54 * time.Second
    pongWait   = 60 * time.Second
    writeWait  = 10 * time.Second
)

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                logger.Error("WebSocket error", "error", err)
            }
            break
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

---

## 1.4 最佳实践

### 1.4.1 WebSocket 设计最佳实践

**为什么需要最佳实践？**

合理的 WebSocket 设计可以提高系统的稳定性、性能和可维护性。

**最佳实践原则**:

1. **连接管理**: 使用 Hub 模式管理连接
2. **消息格式**: 使用统一的消息格式
3. **错误处理**: 完善的错误处理和重连机制
4. **安全性**: 验证 Origin，限制连接数

**实际应用示例**:

```go
// WebSocket 最佳实践
type WebSocketServer struct {
    hub    *Hub
    router *chi.Mux
}

func NewWebSocketServer() *WebSocketServer {
    hub := NewHub()
    go hub.Run()

    return &WebSocketServer{
        hub:    hub,
        router: chi.NewRouter(),
    }
}

// 安全的 WebSocket 升级
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        // 验证 Origin
        return isValidOrigin(origin)
    },
}

// 限制连接数
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    if s.hub.ClientCount() >= maxConnections {
        http.Error(w, "Too many connections", http.StatusServiceUnavailable)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error("Failed to upgrade", "error", err)
        return
    }

    client := &Client{
        hub:  s.hub,
        conn: conn,
        send: make(chan []byte, 256),
    }

    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}
```

**最佳实践要点**:

1. **连接管理**: 使用 Hub 模式集中管理连接
2. **消息格式**: 使用统一的 JSON 消息格式
3. **错误处理**: 完善的错误处理和日志记录
4. **安全性**: 验证 Origin，限制连接数和消息大小

---

## 📚 扩展阅读

- [WebSocket 官方文档](https://github.com/gorilla/websocket)
- [RFC 6455 标准](https://tools.ietf.org/html/rfc6455)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 WebSocket 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
