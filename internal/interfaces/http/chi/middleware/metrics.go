package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Metrics 指标收集器
type Metrics struct {
	mu                sync.RWMutex
	requestCount      map[string]int64      // 请求计数
	requestDuration   map[string]time.Duration // 请求耗时
	errorCount        map[string]int64      // 错误计数
	activeRequests    int64                 // 活跃请求数
	totalRequests     int64                 // 总请求数
	totalDuration     time.Duration         // 总耗时
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

// MetricsMiddleware 性能监控中间件
func MetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			start := time.Now()

			// 增加请求计数
			metrics.IncrementRequest(path)

			// 创建响应包装器
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// 执行下一个处理器
			next.ServeHTTP(ww, r)

			// 记录耗时
			duration := time.Since(start)
			metrics.RecordDuration(path, duration)

			// 减少活跃请求数
			metrics.DecrementActive()

			// 记录错误
			if ww.Status() >= 400 {
				metrics.IncrementError(path)
			}

			// 添加性能指标到响应头
			ww.Header().Set("X-Response-Time", duration.String())
			ww.Header().Set("X-Request-Duration-Ms", strconv.FormatInt(duration.Milliseconds(), 10))
		})
	}
}

// MetricsHandler 指标查询处理器
func MetricsHandler(metrics *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := metrics.GetStats()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// 简单的JSON输出（实际应该使用json包）
		json := `{"total_requests":` + strconv.FormatInt(stats["total_requests"].(int64), 10) +
			`,"active_requests":` + strconv.FormatInt(stats["active_requests"].(int64), 10) +
			`,"average_duration":"` + stats["average_duration"].(string) + `"}`
		w.Write([]byte(json))
	}
}
