package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Metrics 是 HTTP 请求指标收集器。
//
// 功能说明：
// - 收集 HTTP 请求的统计信息
// - 支持按路径统计请求数、耗时和错误数
// - 提供全局统计信息（总请求数、活跃请求数、平均耗时）
//
// 收集的指标：
// - 请求计数：按路径统计请求数
// - 请求耗时：按路径统计请求耗时
// - 错误计数：按路径统计错误数
// - 活跃请求数：当前正在处理的请求数
// - 总请求数：累计处理的请求数
// - 总耗时：累计请求耗时
//
// 字段说明：
// - requestCount: 按路径的请求计数
// - requestDuration: 按路径的请求耗时
// - errorCount: 按路径的错误计数
// - activeRequests: 当前活跃请求数
// - totalRequests: 总请求数
// - totalDuration: 总耗时
// - mu: 读写互斥锁（保证并发安全）
type Metrics struct {
	mu                sync.RWMutex
	requestCount      map[string]int64      // 请求计数
	requestDuration   map[string]time.Duration // 请求耗时
	errorCount        map[string]int64      // 错误计数
	activeRequests    int64                 // 活跃请求数
	totalRequests     int64                 // 总请求数
	totalDuration     time.Duration         // 总耗时
}

// NewMetrics 创建并初始化指标收集器。
//
// 功能说明：
// - 创建新的指标收集器实例
// - 初始化所有指标映射表
//
// 返回：
// - *Metrics: 配置好的指标收集器实例
//
// 使用示例：
//
//	metrics := middleware.NewMetrics()
//	router.Use(middleware.MetricsMiddleware(metrics))
func NewMetrics() *Metrics {
	return &Metrics{
		requestCount:    make(map[string]int64),
		requestDuration: make(map[string]time.Duration),
		errorCount:      make(map[string]int64),
	}
}

// GetStats 获取统计信息。
//
// 功能说明：
// - 返回所有收集的统计信息
// - 计算平均请求耗时
// - 线程安全
//
// 返回：
// - map[string]interface{}: 统计信息
//   - total_requests: 总请求数
//   - active_requests: 活跃请求数
//   - average_duration: 平均请求耗时
//   - request_count: 按路径的请求计数
//   - error_count: 按路径的错误计数
//
// 使用示例：
//
//	stats := metrics.GetStats()
//	fmt.Printf("Total requests: %d\n", stats["total_requests"])
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

// IncrementRequest 增加请求计数。
//
// 功能说明：
// - 增加指定路径的请求计数
// - 增加总请求数
// - 增加活跃请求数
//
// 参数：
// - path: 请求路径
//
// 使用场景：
// - 在请求开始时调用
func (m *Metrics) IncrementRequest(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCount[path]++
	m.totalRequests++
	m.activeRequests++
}

// DecrementActive 减少活跃请求数。
//
// 功能说明：
// - 在请求完成时减少活跃请求数
//
// 使用场景：
// - 在请求结束时调用
func (m *Metrics) DecrementActive() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeRequests--
}

// RecordDuration 记录请求耗时。
//
// 功能说明：
// - 记录指定路径的请求耗时
// - 累加到总耗时中
//
// 参数：
// - path: 请求路径
// - duration: 请求耗时
//
// 使用场景：
// - 在请求完成时调用
func (m *Metrics) RecordDuration(path string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestDuration[path] += duration
	m.totalDuration += duration
}

// IncrementError 增加错误计数。
//
// 功能说明：
// - 增加指定路径的错误计数
//
// 参数：
// - path: 请求路径
//
// 使用场景：
// - 在请求返回错误状态码时调用
func (m *Metrics) IncrementError(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCount[path]++
}

// MetricsMiddleware 创建性能监控中间件。
//
// 功能说明：
// - 收集 HTTP 请求的性能指标
// - 记录请求计数、耗时和错误数
// - 在响应头中添加性能信息
//
// 工作流程：
// 1. 记录请求开始时间
// 2. 增加请求计数和活跃请求数
// 3. 执行下一个处理器
// 4. 计算请求耗时
// 5. 记录耗时和减少活跃请求数
// 6. 如果状态码 >= 400，记录错误
// 7. 在响应头中添加性能信息
//
// 响应头：
// - X-Response-Time: 请求耗时（字符串格式）
// - X-Request-Duration-Ms: 请求耗时（毫秒）
//
// 参数：
// - metrics: 指标收集器实例
//
// 返回：
// - func(http.Handler) http.Handler: Chi 中间件函数
//
// 使用示例：
//
//	metrics := middleware.NewMetrics()
//	router.Use(middleware.MetricsMiddleware(metrics))
//
//	// 查询指标
//	router.Get("/metrics", middleware.MetricsHandler(metrics))
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

// MetricsHandler 创建指标查询处理器。
//
// 功能说明：
// - 提供 HTTP 端点查询指标
// - 返回 JSON 格式的统计信息
// - 可用于监控系统集成（如 Prometheus）
//
// 返回的指标：
// - total_requests: 总请求数
// - active_requests: 活跃请求数
// - average_duration: 平均请求耗时
//
// 参数：
// - metrics: 指标收集器实例
//
// 返回：
// - http.HandlerFunc: HTTP 处理器函数
//
// 使用示例：
//
//	metrics := middleware.NewMetrics()
//	router.Get("/metrics", middleware.MetricsHandler(metrics))
//
// 注意事项：
// - 当前实现使用简单的字符串拼接生成 JSON
// - 生产环境建议使用 json 包进行序列化
// - 可以扩展以支持 Prometheus 格式
func MetricsHandler(metrics *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := metrics.GetStats()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// 简单的JSON输出（实际应该使用json包）
		// TODO: 使用 encoding/json 包进行序列化
		json := `{"total_requests":` + strconv.FormatInt(stats["total_requests"].(int64), 10) +
			`,"active_requests":` + strconv.FormatInt(stats["active_requests"].(int64), 10) +
			`,"average_duration":"` + stats["average_duration"].(string) + `"}`
		w.Write([]byte(json))
	}
}
