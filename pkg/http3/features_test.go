package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// =============================================================================
// 增强功能测试
// =============================================================================

// TestWebSocketHub WebSocket Hub测试
func TestWebSocketHub(t *testing.T) {
	hub := NewWSHub()

	// 测试初始状态
	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", hub.ClientCount())
	}

	// 启动hub
	go hub.Run()

	// 等待hub启动
	time.Sleep(100 * time.Millisecond)

	// 测试广播
	err := hub.Broadcast("test", "hello")
	if err != nil {
		t.Errorf("Failed to broadcast: %v", err)
	}
}

// TestMiddlewareChain 中间件链测试
func TestMiddlewareChain(t *testing.T) {
	called := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	chain := NewMiddlewareChain()
	chain.Use(LoggingMiddleware)
	chain.Use(RecoveryMiddleware)

	finalHandler := chain.Then(handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	finalHandler.ServeHTTP(w, req)

	if !called {
		t.Error("Handler was not called")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestLoggingMiddleware 日志中间件测试
func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestRecoveryMiddleware 恢复中间件测试
func TestRecoveryMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := RecoveryMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 不应该panic
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// TestCORSMiddleware CORS中间件测试
func TestCORSMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CORSMiddleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// 检查CORS头
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS headers not set")
	}
}

// TestTimeoutMiddleware 超时中间件测试
func TestTimeoutMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	middleware := TimeoutMiddleware(100 * time.Millisecond)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// 应该超时
	if w.Code != http.StatusRequestTimeout {
		t.Errorf("Expected timeout status 408, got %d", w.Code)
	}
}

// TestRequestIDMiddleware 请求ID中间件测试
func TestRequestIDMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value("request_id")
		if requestID == nil {
			t.Error("Request ID not set in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header not set")
	}
}

// TestSecurityHeadersMiddleware 安全头中间件测试
func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := SecurityHeadersMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// 检查安全头
	headers := []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Strict-Transport-Security",
		"Content-Security-Policy",
	}

	for _, header := range headers {
		if w.Header().Get(header) == "" {
			t.Errorf("Security header %s not set", header)
		}
	}
}

// TestConnectionPool 连接池测试
func TestConnectionPool(t *testing.T) {
	pool := NewConnectionPool(5)

	// 测试获取连接
	client, err := pool.Get()
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	if client == nil {
		t.Fatal("Got nil client")
	}

	// 测试归还连接
	pool.Put(client)

	// 测试统计
	stats := pool.Stats()
	if stats["max_connections"] != 5 {
		t.Errorf("Expected max_connections 5, got %v", stats["max_connections"])
	}

	// 清理
	pool.Close()
}

// TestConnectionManager 连接管理器测试
func TestConnectionManager(t *testing.T) {
	manager := NewConnectionManager(100)

	// 跟踪连接
	conn := manager.Track("conn1", "192.168.1.1:8080")
	if conn == nil {
		t.Fatal("Failed to track connection")
	}

	// 更新连接
	manager.Update("conn1", 100, 200)

	// 获取统计
	stats := manager.GetStats()
	if stats["active_connections"] != 1 {
		t.Errorf("Expected 1 active connection, got %v", stats["active_connections"])
	}

	// 移除连接
	manager.Remove("conn1")

	stats = manager.GetStats()
	if stats["active_connections"] != 0 {
		t.Errorf("Expected 0 active connections after removal, got %v", stats["active_connections"])
	}
}

// TestConnectionCleanup 连接清理测试
func TestConnectionCleanup(t *testing.T) {
	manager := NewConnectionManager(100)

	// 添加一些连接
	manager.Track("conn1", "192.168.1.1:8080")
	manager.Track("conn2", "192.168.1.2:8080")

	// 等待一段时间
	time.Sleep(100 * time.Millisecond)

	// 清理空闲连接
	cleaned := manager.Cleanup(50 * time.Millisecond)

	if cleaned != 2 {
		t.Errorf("Expected 2 connections cleaned, got %d", cleaned)
	}
}

// TestCacheMiddleware 缓存中间件测试
func TestCacheMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CacheMiddleware(1 * time.Hour)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	if cacheControl == "" {
		t.Error("Cache-Control header not set")
	}
}

// =============================================================================
// 基准测试
// =============================================================================

// BenchmarkMiddlewareChain 中间件链基准测试
func BenchmarkMiddlewareChain(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chain := NewMiddlewareChain()
	chain.Use(LoggingMiddleware)
	chain.Use(RecoveryMiddleware)
	chain.Use(CORSMiddleware)

	finalHandler := chain.Then(handler)

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		finalHandler.ServeHTTP(w, req)
	}
}

// BenchmarkConnectionPool 连接池基准测试
func BenchmarkConnectionPool(b *testing.B) {
	pool := NewConnectionPool(100)
	defer pool.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := pool.Get()
		if err != nil {
			b.Fatal(err)
		}
		pool.Put(client)
	}
}

// BenchmarkConnectionManager 连接管理器基准测试
func BenchmarkConnectionManager(b *testing.B) {
	manager := NewConnectionManager(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := generateClientID()
		manager.Track(id, "192.168.1.1:8080")
		manager.Update(id, 100, 200)
		manager.Remove(id)
	}
}

// BenchmarkSecurityHeaders 安全头基准测试
func BenchmarkSecurityHeaders(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := SecurityHeadersMiddleware(handler)
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
	}
}
