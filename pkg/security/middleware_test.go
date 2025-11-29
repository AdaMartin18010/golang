package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSecurityMiddleware_Middleware(t *testing.T) {
	config := SecurityMiddlewareConfig{
		SecurityHeaders: func() *SecurityHeadersConfig {
			cfg := DefaultSecurityHeadersConfig()
			return &cfg
		}(),
		RateLimit: &RateLimiterConfig{
			Limit:  100,
			Window: 1 * time.Minute,
		},
	}

	middleware := NewSecurityMiddleware(config)
	defer middleware.Shutdown()

	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestSecurityMiddleware_RateLimit(t *testing.T) {
	config := SecurityMiddlewareConfig{
		RateLimit: &RateLimiterConfig{
			Limit:  1,
			Window: 1 * time.Minute,
		},
	}

	middleware := NewSecurityMiddleware(config)
	defer middleware.Shutdown()

	handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// 第一次请求应该成功
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("First request should succeed, got %d", rr.Code)
	}

	// 第二次请求应该被限制
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request should be rate limited, got %d", rr2.Code)
	}
}

func TestSecurityMiddleware_CSRF(t *testing.T) {
	csrfConfig := DefaultCSRFConfig()
	config := SecurityMiddlewareConfig{
		CSRF: &csrfConfig,
	}

	middleware := NewSecurityMiddleware(config)
	defer middleware.Shutdown()

	// 创建会话和令牌
	sessionID := "test-session"
	token, _ := middleware.csrfProtection.GenerateToken(sessionID)

	handler := middleware.CSRFMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 测试 POST 请求（需要 CSRF 令牌）
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("X-CSRF-Token", token)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Request with valid CSRF token should succeed, got %d", rr.Code)
	}

	// 测试无效令牌
	req2 := httptest.NewRequest("POST", "/", nil)
	req2.Header.Set("X-CSRF-Token", "invalid-token")
	req2.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusForbidden {
		t.Errorf("Request with invalid CSRF token should fail, got %d", rr2.Code)
	}
}

func TestSecurityMiddleware_SecurityHeaders(t *testing.T) {
	headersConfig := DefaultSecurityHeadersConfig()
	config := SecurityMiddlewareConfig{
		SecurityHeaders: &headersConfig,
	}

	middleware := NewSecurityMiddleware(config)
	defer middleware.Shutdown()

	handler := middleware.SecurityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// 检查安全头部
	if rr.Header().Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header should be set")
	}

	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("X-Content-Type-Options header should be set")
	}
}
