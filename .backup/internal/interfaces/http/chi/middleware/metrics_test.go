package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestMetricsMiddleware(t *testing.T) {
	metrics := NewMetrics()
	r := chi.NewRouter()
	r.Use(MetricsMiddleware(metrics))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应头
	if w.Header().Get("X-Response-Time") == "" {
		t.Error("Expected X-Response-Time header")
	}

	// 检查统计信息
	stats := metrics.GetStats()
	if stats["total_requests"].(int64) != 1 {
		t.Errorf("Expected 1 total request, got %d", stats["total_requests"].(int64))
	}
}

func TestMetrics_IncrementError(t *testing.T) {
	metrics := NewMetrics()
	metrics.IncrementRequest("/test")
	metrics.IncrementError("/test")

	stats := metrics.GetStats()
	errorCount := stats["error_count"].(map[string]int64)
	if errorCount["/test"] != 1 {
		t.Errorf("Expected 1 error, got %d", errorCount["/test"])
	}
}

func TestMetricsHandler(t *testing.T) {
	metrics := NewMetrics()
	metrics.IncrementRequest("/test")

	handler := MetricsHandler(metrics)
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}
