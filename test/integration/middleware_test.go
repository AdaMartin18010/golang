package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/golang/pkg/observability/operational"
	"github.com/yourusername/golang/pkg/security"
)

// TestMiddlewareIntegration 测试中间件集成
func TestMiddlewareIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建测试处理器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 测试日志中间件
	loggingMiddleware := operational.LoggingMiddleware()
	loggedHandler := loggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

// TestRecoveryMiddlewareIntegration 测试恢复中间件
func TestRecoveryMiddlewareIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建会panic的处理器
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// 使用恢复中间件
	recoveryMiddleware := operational.RecoveryMiddleware()
	recoveredHandler := recoveryMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// 应该不会panic，而是返回500错误
	recoveredHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}
}

// TestCORSMiddlewareIntegration 测试CORS中间件
func TestCORSMiddlewareIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建CORS配置
	corsConfig := operational.CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         3600,
	}

	// 创建CORS中间件
	corsMiddleware := operational.CORSMiddleware(corsConfig)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := corsMiddleware(handler)

	// 测试预检请求
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rr := httptest.NewRecorder()

	corsHandler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("CORS headers should be set")
	}
}

// TestSecurityHeadersMiddlewareIntegration 测试安全头部中间件
func TestSecurityHeadersMiddlewareIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建安全头部配置
	headersConfig := security.DefaultSecurityHeadersConfig()
	headers := security.NewSecurityHeaders(headersConfig)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	headersHandler := headers.Middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	headersHandler.ServeHTTP(rr, req)

	// 验证安全头部
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("X-Content-Type-Options header should be set")
	}

	if rr.Header().Get("X-Frame-Options") == "" {
		t.Error("X-Frame-Options header should be set")
	}
}

// TestRateLimitMiddlewareIntegration 测试速率限制中间件
func TestRateLimitMiddlewareIntegration(t *testing.T) {
	config := DefaultTestFrameworkConfig()
	tf, err := NewTestFramework(config)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	defer tf.Shutdown()

	// 创建速率限制器
	limiter := security.NewIPRateLimiter(security.RateLimiterConfig{
		Limit:  2,
		Window: 1 * time.Minute,
	})
	defer limiter.Shutdown(context.Background())

	// 创建速率限制中间件
	rateLimitMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			allowed, err := limiter.AllowIP(r.Context(), ip)
			if err != nil || !allowed {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limitedHandler := rateLimitMiddleware(handler)

	// 前两次请求应该成功
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1"
		rr := httptest.NewRecorder()

		limitedHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Request %d should succeed, got %d", i+1, rr.Code)
		}
	}

	// 第三次请求应该被限制
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1"
	rr := httptest.NewRecorder()

	limitedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request should be rate limited, got %d", rr.Code)
	}
}
