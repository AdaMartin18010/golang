package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// =============================================================================
// WebSocket支持 - WebSocket Support
// =============================================================================

// WSMessage WebSocket消息结构
type WSMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// WSClient WebSocket客户端
type WSClient struct {
	conn      *websocket.Conn
	send      chan []byte
	hub       *WSHub
	id        string
	createdAt time.Time
}

// WSHub WebSocket管理中心
type WSHub struct {
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	mu         sync.RWMutex
}

// upgrader WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应该检查Origin
		return true
	},
}

// NewWSHub 创建WebSocket中心
func NewWSHub() *WSHub {
	return &WSHub{
		clients:    make(map[*WSClient]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

// Run 运行WebSocket中心
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client registered: %s", client.id)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("WebSocket client unregistered: %s", client.id)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 客户端发送缓冲区满，关闭连接
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// ClientCount 获取连接客户端数量
func (h *WSHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Broadcast 广播消息到所有客户端
func (h *WSHub) Broadcast(msgType string, data interface{}) error {
	msg := WSMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.broadcast <- jsonData
	return nil
}

// readPump 读取WebSocket消息
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 回显消息（示例）
		response := WSMessage{
			Type:      "echo",
			Data:      msg.Data,
			Timestamp: time.Now(),
		}

		jsonData, _ := json.Marshal(response)
		c.send <- jsonData
	}
}

// writePump 发送WebSocket消息
func (c *WSClient) writePump() {
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
				// Hub关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送队列中的消息
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

// handleWebSocket WebSocket处理器
func handleWebSocket(hub *WSHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		client := &WSClient{
			conn:      conn,
			send:      make(chan []byte, 256),
			hub:       hub,
			id:        generateClientID(),
			createdAt: time.Now(),
		}

		client.hub.register <- client

		// 在新的goroutine中运行读写pump
		go client.writePump()
		go client.readPump()
	}
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}

// handleWSStats WebSocket统计信息
func handleWSStats(hub *WSHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := map[string]interface{}{
			"connected_clients": hub.ClientCount(),
			"timestamp":         time.Now(),
			"status":            "active",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
