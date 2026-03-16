package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestCircuitBreaker_Allow(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		Timeout:          100 * time.Millisecond,
		TimeoutWindow:    1 * time.Second,
	}

	breaker := NewCircuitBreaker("test", config)

	// 初始状态应该是关闭的
	if !breaker.Allow() {
		t.Error("Expected circuit breaker to allow requests in closed state")
	}

	// 触发失败
	for i := 0; i < 3; i++ {
		breaker.OnFailure()
	}

	// 现在应该被熔断
	if breaker.Allow() {
		t.Error("Expected circuit breaker to block requests in open state")
	}

	// 等待超时后应该进入半开状态
	time.Sleep(150 * time.Millisecond)
	if !breaker.Allow() {
		t.Error("Expected circuit breaker to allow requests in half-open state")
	}

	// 在半开状态下成功2次应该关闭熔断器
	breaker.OnSuccess()
	breaker.OnSuccess()

	if breaker.GetState() != StateClosed {
		t.Error("Expected circuit breaker to be closed after successful recovery")
	}
}

func TestCircuitBreakerMiddleware(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}

	breaker := NewCircuitBreaker("test", config)

	r := chi.NewRouter()
	r.Use(CircuitBreakerMiddleware(breaker))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	// 发送失败请求
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	// 现在应该被熔断
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", w.Code)
	}
}

func TestCircuitBreakerManager(t *testing.T) {
	manager := NewCircuitBreakerManager()

	config := CircuitBreakerConfig{
		FailureThreshold: 5,
		Timeout:          time.Second,
	}

	// 获取或创建熔断器
	breaker1 := manager.GetOrCreate("test1", config)
	breaker2 := manager.GetOrCreate("test1", config)

	// 应该是同一个实例
	if breaker1 != breaker2 {
		t.Error("Expected same circuit breaker instance")
	}

	// 获取不同的熔断器
	breaker3 := manager.GetOrCreate("test2", config)
	if breaker1 == breaker3 {
		t.Error("Expected different circuit breaker instances")
	}

	// 获取不存在的熔断器
	_, exists := manager.Get("nonexistent")
	if exists {
		t.Error("Expected circuit breaker to not exist")
	}
}
