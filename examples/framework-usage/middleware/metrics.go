package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Metrics 是HTTP请求指标收集器
type Metrics struct {
	mu              sync.RWMutex
	requestCount    map[string]int64
	requestDuration map[string]time.Duration
	errorCount      map[string]int64
	activeRequests  int64
	totalRequests   int64
	totalDuration   time.Duration
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		requestCount:    make(map[string]int64),
		requestDuration: make(map[string]time.Duration),
		errorCount:      make(map[string]int64),
	}
}

// GetStats 获取统计信息
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	avgDuration := time.Duration(0)
	if m.totalRequests > 0 {
		avgDuration = m.totalDuration / time.Duration(m.totalRequests)
	}

	return map[string]interface{}{
		"total_requests":   m.totalRequests,
		"active_requests":  m.activeRequests,
		"average_duration": avgDuration.String(),
		"request_count":    m.requestCount,
		"error_count":      m.errorCount,
	}
}

// IncrementRequest 增加请求计数
func (m *Metrics) IncrementRequest(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount[path]++
	m.totalRequests++
	m.activeRequests++
}

// DecrementActive 减少活跃请求数
func (m *Metrics) DecrementActive() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeRequests--
}

// RecordDuration 记录请求耗时
func (m *Metrics) RecordDuration(path string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestDuration[path] += duration
	m.totalDuration += duration
}

// IncrementError 增加错误计数
func (m *Metrics) IncrementError(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCount[path]++
}

// responseWriter 是包装http.ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware 创建性能监控中间件
func MetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			start := time.Now()

			metrics.IncrementRequest(path)
			defer metrics.DecrementActive()

			wrapped := newResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			metrics.RecordDuration(path, duration)

			if wrapped.statusCode >= 400 {
				metrics.IncrementError(path)
			}

			// 添加性能指标到响应头
			w.Header().Set("X-Response-Time", duration.String())
			w.Header().Set("X-Request-Duration-Ms", strconv.FormatInt(duration.Milliseconds(), 10))
		})
	}
}

// MetricsHandler 创建指标查询处理器
func MetricsHandler(metrics *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := metrics.GetStats()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}
