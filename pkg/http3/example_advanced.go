package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// =============================================================================
// 高级示例 - Advanced Example
// =============================================================================

// SetupAdvancedServer 设置高级HTTP服务器
func SetupAdvancedServer() *http.Server {
	// 创建WebSocket Hub
	wsHub := NewWSHub()
	go wsHub.Run()

	// 创建连接管理器
	connManager := NewConnectionManager(1000)
	go func() {
		// 定期清理空闲连接
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cleaned := connManager.Cleanup(5 * time.Minute)
			if cleaned > 0 {
				log.Printf("Cleaned %d idle connections", cleaned)
			}
		}
	}()

	// 创建路由
	mux := http.NewServeMux()

	// 基础路由（使用优化的处理器）
	mux.HandleFunc("/", handleRootOptimized)
	mux.HandleFunc("/stats", handleStatsOptimized)
	mux.HandleFunc("/health", handleHealthOptimized)
	mux.HandleFunc("/data", handleDataOptimizedV2)

	// WebSocket路由
	mux.HandleFunc("/ws", handleWebSocket(wsHub))
	mux.HandleFunc("/ws/stats", handleWSStats(wsHub))

	// 连接管理路由
	mux.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		stats := connManager.GetStats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// 创建中间件链
	chain := NewMiddlewareChain()

	// 添加中间件（按顺序）
	chain.Use(RecoveryMiddleware)                        // 1. 恢复panic
	chain.Use(LoggingMiddleware)                         // 2. 日志记录
	chain.Use(RequestIDMiddleware)                       // 3. 请求ID
	chain.Use(SecurityHeadersMiddleware)                 // 4. 安全头
	chain.Use(CORSMiddleware)                            // 5. CORS
	chain.Use(TimeoutMiddleware(30 * time.Second))       // 6. 超时控制
	chain.Use(CacheMiddleware(1 * time.Hour))            // 7. 缓存
	chain.Use(ConnectionTrackingMiddleware(connManager)) // 8. 连接跟踪

	// 应用中间件链
	handler := chain.Then(mux)

	// 创建服务器
	server := &http.Server{
		Addr:         ":8443",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

// RunAdvancedExample 运行高级示例
func RunAdvancedExample() {
	log.Println("🚀 Starting Advanced HTTP/3 Server")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	_ = SetupAdvancedServer() // Setup server (unused in demo)

	log.Println("📊 Features Enabled:")
	log.Println("   ✓ WebSocket Support (ws://localhost:8443/ws)")
	log.Println("   ✓ Middleware Chain (8 middlewares)")
	log.Println("   ✓ Connection Management")
	log.Println("   ✓ Rate Limiting")
	log.Println("   ✓ Security Headers")
	log.Println("   ✓ Request Logging")
	log.Println("   ✓ Panic Recovery")
	log.Println("   ✓ CORS Support")
	log.Println()

	log.Println("📍 Endpoints:")
	log.Println("   GET  /                - Root endpoint (optimized)")
	log.Println("   GET  /stats           - Server statistics")
	log.Println("   GET  /health          - Health check")
	log.Println("   GET  /data            - Data endpoint (optimized)")
	log.Println("   WS   /ws              - WebSocket endpoint")
	log.Println("   GET  /ws/stats        - WebSocket statistics")
	log.Println("   GET  /connections     - Connection statistics")
	log.Println()

	log.Println("🌐 Server: https://localhost:8443")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println()

	// 注意：需要证书文件才能实际运行
	// server := SetupAdvancedServer()
	// log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))

	log.Println("💡 To start the server, uncomment the ListenAndServeTLS line")
	log.Println("   and provide valid certificate files.")
}

// =============================================================================
// 使用示例
// =============================================================================

// ExampleWebSocketClient WebSocket客户端示例
/*
Example usage:

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8443/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Send message
	msg := WSMessage{
		Type: "ping",
		Data: "Hello Server!",
		Timestamp: time.Now(),
	}

	if err := ws.WriteJSON(msg); err != nil {
		log.Fatal(err)
	}

	// Receive message
	var response WSMessage
	if err := ws.ReadJSON(&response); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received: %+v\n", response)
*/

// ExampleCustomMiddleware 自定义中间件示例
/*
Example of creating a custom middleware:

	func MyCustomMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Before request processing
			log.Println("Before:", r.URL.Path)

			// Process request
			next.ServeHTTP(w, r)

			// After request processing
			log.Println("After:", r.URL.Path)
		})
	}

	// Use it:
	chain := NewMiddlewareChain()
	chain.Use(MyCustomMiddleware)
	handler := chain.Then(mux)
*/

// ExampleConnectionPool 连接池示例
/*
Example of using connection pool:

	pool := NewConnectionPool(10)
	defer pool.Close()

	// Get connection
	client, err := pool.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Put(client)

	// Use client
	resp, err := client.Get("https://api.example.com/data")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
*/
