package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestTracingMiddleware(t *testing.T) {
	config := TracingConfig{
		TracerName:     "test-tracer",
		ServiceName:    "test-service",
		ServiceVersion: "v1.0.0",
		SkipPaths:      []string{"/health"},
	}

	r := chi.NewRouter()
	r.Use(TracingMiddleware(config))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应头
	if w.Header().Get("X-Trace-ID") == "" {
		t.Error("Expected X-Trace-ID header")
	}
	if w.Header().Get("X-Span-ID") == "" {
		t.Error("Expected X-Span-ID header")
	}
}

func TestTracingMiddleware_SkipPaths(t *testing.T) {
	config := TracingConfig{
		SkipPaths: []string{"/health"},
	}

	r := chi.NewRouter()
	r.Use(TracingMiddleware(config))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 跳过的路径不应该有追踪头
	if w.Header().Get("X-Trace-ID") != "" {
		t.Error("Expected no X-Trace-ID header for skipped path")
	}
}
