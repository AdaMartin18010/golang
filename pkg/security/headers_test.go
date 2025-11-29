package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders_Middleware(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	headers := NewSecurityHeaders(config)

	handler := headers.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	if rr.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("X-Frame-Options header should be set")
	}

	if rr.Header().Get("X-XSS-Protection") == "" {
		t.Error("X-XSS-Protection header should be set")
	}

	// 检查服务器信息已移除
	if rr.Header().Get("Server") != "" {
		t.Error("Server header should be removed")
	}

	if rr.Header().Get("X-Powered-By") != "" {
		t.Error("X-Powered-By header should be removed")
	}
}

func TestSecurityHeaders_HandlerFunc(t *testing.T) {
	config := DefaultSecurityHeadersConfig()
	headers := NewSecurityHeaders(config)

	handler := headers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	// 检查安全头部
	if rr.Header().Get("Content-Security-Policy") == "" {
		t.Error("Content-Security-Policy header should be set")
	}
}
