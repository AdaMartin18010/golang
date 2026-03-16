package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestRateLimitMiddleware(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 2,
		Burst:             2,
		Window:            time.Second,
	}

	r := chi.NewRouter()
	r.Use(RateLimitMiddleware(config))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 发送多个请求
	allowed := 0
	blocked := 0

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			allowed++
		} else if w.Code == http.StatusTooManyRequests {
			blocked++
		}
	}

	// 应该允许前2个请求，阻止后续请求
	if allowed < 2 {
		t.Errorf("Expected at least 2 allowed requests, got %d", allowed)
	}
	if blocked == 0 {
		t.Error("Expected some blocked requests, got 0")
	}
}

func TestTokenBucket_Allow(t *testing.T) {
	bucket := NewTokenBucket(2, 1.0) // 容量2，每秒填充1个令牌

	// 前2个请求应该允许
	if !bucket.Allow() {
		t.Error("Expected first request to be allowed")
	}
	if !bucket.Allow() {
		t.Error("Expected second request to be allowed")
	}

	// 第3个请求应该被阻止
	if bucket.Allow() {
		t.Error("Expected third request to be blocked")
	}

	// 等待1秒后应该允许
	time.Sleep(1100 * time.Millisecond)
	if !bucket.Allow() {
		t.Error("Expected request after refill to be allowed")
	}
}

func TestRateLimitMiddleware_SkipPaths(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 1,
		Burst:             1,
		SkipPaths:         []string{"/public"},
	}

	r := chi.NewRouter()
	r.Use(RateLimitMiddleware(config))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/public", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 发送多个请求到跳过路径
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/public", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for skipped path, got %d", w.Code)
		}
	}
}
