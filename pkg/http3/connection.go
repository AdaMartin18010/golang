package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// =============================================================================
// 连接管理 - Connection Management
// =============================================================================

// ConnectionPool HTTP连接池
type ConnectionPool struct {
	maxConnections int
	connections    chan *http.Client
	mu             sync.RWMutex
	activeCount    int
	totalRequests  int64
	createdAt      time.Time
}

// NewConnectionPool 创建连接池
func NewConnectionPool(maxConn int) *ConnectionPool {
	pool := &ConnectionPool{
		maxConnections: maxConn,
		connections:    make(chan *http.Client, maxConn),
		createdAt:      time.Now(),
	}

	// 预创建连接
	for i := 0; i < maxConn; i++ {
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		}
		pool.connections <- client
	}

	return pool
}

// Get 获取连接
func (p *ConnectionPool) Get() (*http.Client, error) {
	p.mu.Lock()
	p.activeCount++
	p.totalRequests++
	p.mu.Unlock()

	select {
	case client := <-p.connections:
		return client, nil
	case <-time.After(5 * time.Second):
		p.mu.Lock()
		p.activeCount--
		p.mu.Unlock()
		return nil, fmt.Errorf("connection pool timeout")
	}
}

// Put 归还连接
func (p *ConnectionPool) Put(client *http.Client) {
	p.mu.Lock()
	p.activeCount--
	p.mu.Unlock()

	select {
	case p.connections <- client:
	default:
		// 连接池已满，关闭客户端
		client.CloseIdleConnections()
	}
}

// Stats 获取连接池统计
func (p *ConnectionPool) Stats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"max_connections": p.maxConnections,
		"active_count":    p.activeCount,
		"idle_count":      len(p.connections),
		"total_requests":  p.totalRequests,
		"uptime":          time.Since(p.createdAt).String(),
	}
}

// Close 关闭连接池
func (p *ConnectionPool) Close() {
	close(p.connections)
	for client := range p.connections {
		client.CloseIdleConnections()
	}
}

// =============================================================================
// 连接管理器
// =============================================================================

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*Connection
	mu          sync.RWMutex
	maxConn     int
}

// Connection 连接信息
type Connection struct {
	ID         string
	RemoteAddr string
	CreatedAt  time.Time
	LastActive time.Time
	Requests   int64
	BytesSent  int64
	BytesRecv  int64
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(maxConn int) *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
		maxConn:     maxConn,
	}
}

// Track 跟踪连接
func (cm *ConnectionManager) Track(connID, remoteAddr string) *Connection {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn := &Connection{
		ID:         connID,
		RemoteAddr: remoteAddr,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}

	cm.connections[connID] = conn
	return conn
}

// Update 更新连接信息
func (cm *ConnectionManager) Update(connID string, bytesSent, bytesRecv int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn, ok := cm.connections[connID]; ok {
		conn.LastActive = time.Now()
		conn.Requests++
		conn.BytesSent += bytesSent
		conn.BytesRecv += bytesRecv
	}
}

// Remove 移除连接
func (cm *ConnectionManager) Remove(connID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.connections, connID)
}

// GetStats 获取统计信息
func (cm *ConnectionManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var totalRequests, totalBytesSent, totalBytesRecv int64
	for _, conn := range cm.connections {
		totalRequests += conn.Requests
		totalBytesSent += conn.BytesSent
		totalBytesRecv += conn.BytesRecv
	}

	return map[string]interface{}{
		"active_connections": len(cm.connections),
		"total_requests":     totalRequests,
		"total_bytes_sent":   totalBytesSent,
		"total_bytes_recv":   totalBytesRecv,
	}
}

// Cleanup 清理空闲连接
func (cm *ConnectionManager) Cleanup(idleTimeout time.Duration) int {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cleaned := 0
	now := time.Now()

	for id, conn := range cm.connections {
		if now.Sub(conn.LastActive) > idleTimeout {
			delete(cm.connections, id)
			cleaned++
		}
	}

	return cleaned
}

// =============================================================================
// 连接中间件
// =============================================================================

// ConnectionTrackingMiddleware 连接跟踪中间件
func ConnectionTrackingMiddleware(manager *ConnectionManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			connID := generateClientID()
			ctx := context.WithValue(r.Context(), "connection_id", connID)
			r = r.WithContext(ctx)

			// 跟踪连接
			conn := manager.Track(connID, r.RemoteAddr)
			defer manager.Remove(connID)

			// 创建响应记录器以跟踪字节数
			rec := &bytesRecorder{
				ResponseWriter: w,
				connection:     conn,
				manager:        manager,
			}

			next.ServeHTTP(rec, r)
		})
	}
}

// bytesRecorder 字节记录器
type bytesRecorder struct {
	http.ResponseWriter
	connection *Connection
	manager    *ConnectionManager
	written    int64
}

func (br *bytesRecorder) Write(b []byte) (int, error) {
	n, err := br.ResponseWriter.Write(b)
	br.written += int64(n)

	if br.connection != nil && br.manager != nil {
		br.manager.Update(br.connection.ID, br.written, 0)
	}

	return n, err
}
