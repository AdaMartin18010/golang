package observability

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// =============================================================================
// 指标收集 - Metrics Collection
// =============================================================================

// MetricType 指标类型
type MetricType int

const (
	CounterType MetricType = iota
	GaugeType
	HistogramType
	SummaryType
)

// Metric 指标接口
type Metric interface {
	Name() string
	Type() MetricType
	Value() interface{}
	Labels() map[string]string
}

// =============================================================================
// Counter - 计数器
// =============================================================================

// Counter 计数器（只增不减）
type Counter struct {
	name   string
	help   string
	value  uint64
	labels map[string]string
}

// NewCounter 创建计数器
func NewCounter(name, help string, labels map[string]string) *Counter {
	if labels == nil {
		labels = make(map[string]string)
	}
	return &Counter{
		name:   name,
		help:   help,
		labels: labels,
	}
}

func (c *Counter) Name() string {
	return c.name
}

func (c *Counter) Type() MetricType {
	return CounterType
}

func (c *Counter) Value() interface{} {
	return atomic.LoadUint64(&c.value)
}

func (c *Counter) Labels() map[string]string {
	return c.labels
}

// Inc 增加1
func (c *Counter) Inc() {
	c.Add(1)
}

// Add 增加指定值
func (c *Counter) Add(delta uint64) {
	atomic.AddUint64(&c.value, delta)
}

// Get 获取当前值
func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.value)
}

// =============================================================================
// Gauge - 仪表
// =============================================================================

// Gauge 仪表（可增可减）
type Gauge struct {
	name   string
	help   string
	value  int64
	labels map[string]string
}

// NewGauge 创建仪表
func NewGauge(name, help string, labels map[string]string) *Gauge {
	if labels == nil {
		labels = make(map[string]string)
	}
	return &Gauge{
		name:   name,
		help:   help,
		labels: labels,
	}
}

func (g *Gauge) Name() string {
	return g.name
}

func (g *Gauge) Type() MetricType {
	return GaugeType
}

func (g *Gauge) Value() interface{} {
	return atomic.LoadInt64(&g.value)
}

func (g *Gauge) Labels() map[string]string {
	return g.labels
}

// Set 设置值
func (g *Gauge) Set(value int64) {
	atomic.StoreInt64(&g.value, value)
}

// Inc 增加1
func (g *Gauge) Inc() {
	g.Add(1)
}

// Dec 减少1
func (g *Gauge) Dec() {
	g.Add(-1)
}

// Add 增加指定值
func (g *Gauge) Add(delta int64) {
	atomic.AddInt64(&g.value, delta)
}

// Get 获取当前值
func (g *Gauge) Get() int64 {
	return atomic.LoadInt64(&g.value)
}

// =============================================================================
// Histogram - 直方图
// =============================================================================

// Histogram 直方图（分布统计）
type Histogram struct {
	name    string
	help    string
	buckets []float64
	counts  []uint64
	sum     uint64
	count   uint64
	labels  map[string]string
	mu      sync.RWMutex
}

// NewHistogram 创建直方图
func NewHistogram(name, help string, buckets []float64, labels map[string]string) *Histogram {
	if labels == nil {
		labels = make(map[string]string)
	}
	if buckets == nil {
		// 默认buckets
		buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
	}

	return &Histogram{
		name:    name,
		help:    help,
		buckets: buckets,
		counts:  make([]uint64, len(buckets)+1),
		labels:  labels,
	}
}

func (h *Histogram) Name() string {
	return h.name
}

func (h *Histogram) Type() MetricType {
	return HistogramType
}

func (h *Histogram) Value() interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]interface{}{
		"buckets": h.buckets,
		"counts":  h.counts,
		"sum":     atomic.LoadUint64(&h.sum),
		"count":   atomic.LoadUint64(&h.count),
	}
}

func (h *Histogram) Labels() map[string]string {
	return h.labels
}

// Observe 观察一个值
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 找到合适的bucket
	for i, bucket := range h.buckets {
		if value <= bucket {
			atomic.AddUint64(&h.counts[i], 1)
			break
		}
	}

	// 加入到最后一个bucket（+Inf）
	if value > h.buckets[len(h.buckets)-1] {
		atomic.AddUint64(&h.counts[len(h.counts)-1], 1)
	}

	atomic.AddUint64(&h.sum, uint64(value*1000000)) // 保留微秒精度
	atomic.AddUint64(&h.count, 1)
}

// =============================================================================
// MetricsRegistry - 指标注册表
// =============================================================================

// MetricsRegistry 指标注册表
type MetricsRegistry struct {
	metrics map[string]Metric
	mu      sync.RWMutex
}

// NewMetricsRegistry 创建指标注册表
func NewMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{
		metrics: make(map[string]Metric),
	}
}

// Register 注册指标
func (r *MetricsRegistry) Register(metric Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := metric.Name()
	if _, exists := r.metrics[name]; exists {
		return fmt.Errorf("metric %s already registered", name)
	}

	r.metrics[name] = metric
	return nil
}

// Unregister 注销指标
func (r *MetricsRegistry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.metrics, name)
}

// Get 获取指标
func (r *MetricsRegistry) Get(name string) (Metric, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metric, ok := r.metrics[name]
	return metric, ok
}

// All 获取所有指标
func (r *MetricsRegistry) All() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := make([]Metric, 0, len(r.metrics))
	for _, m := range r.metrics {
		metrics = append(metrics, m)
	}
	return metrics
}

// Export 导出所有指标（Prometheus格式）
func (r *MetricsRegistry) Export() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var output string
	for _, metric := range r.metrics {
		// 添加HELP
		output += fmt.Sprintf("# HELP %s\n", metric.Name())

		// 添加TYPE
		typeStr := "untyped"
		switch metric.Type() {
		case CounterType:
			typeStr = "counter"
		case GaugeType:
			typeStr = "gauge"
		case HistogramType:
			typeStr = "histogram"
		}
		output += fmt.Sprintf("# TYPE %s %s\n", metric.Name(), typeStr)

		// 添加值
		labels := metric.Labels()
		labelStr := ""
		if len(labels) > 0 {
			labelStr = "{"
			first := true
			for k, v := range labels {
				if !first {
					labelStr += ","
				}
				labelStr += fmt.Sprintf("%s=\"%s\"", k, v)
				first = false
			}
			labelStr += "}"
		}

		output += fmt.Sprintf("%s%s %v\n", metric.Name(), labelStr, metric.Value())
		output += "\n"
	}

	return output
}

// =============================================================================
// 全局注册表
// =============================================================================

var (
	defaultRegistry = NewMetricsRegistry()
)

// Register 注册到默认注册表
func Register(metric Metric) error {
	return defaultRegistry.Register(metric)
}

// Unregister 从默认注册表注销
func Unregister(name string) {
	defaultRegistry.Unregister(name)
}

// GetMetric 从默认注册表获取指标
func GetMetric(name string) (Metric, bool) {
	return defaultRegistry.Get(name)
}

// AllMetrics 获取所有指标
func AllMetrics() []Metric {
	return defaultRegistry.All()
}

// ExportMetrics 导出所有指标
func ExportMetrics() string {
	return defaultRegistry.Export()
}

// =============================================================================
// 便捷函数
// =============================================================================

// RegisterCounter 创建并注册计数器
func RegisterCounter(name, help string, labels map[string]string) *Counter {
	counter := NewCounter(name, help, labels)
	// 忽略注册错误，因为重复注册不是致命错误
	_ = Register(counter) // #nosec G104
	return counter
}

// RegisterGauge 创建并注册仪表
func RegisterGauge(name, help string, labels map[string]string) *Gauge {
	gauge := NewGauge(name, help, labels)
	// 忽略注册错误，因为重复注册不是致命错误
	_ = Register(gauge) // #nosec G104
	return gauge
}

// RegisterHistogram 创建并注册直方图
func RegisterHistogram(name, help string, buckets []float64, labels map[string]string) *Histogram {
	histogram := NewHistogram(name, help, buckets, labels)
	// 忽略注册错误，因为重复注册不是致命错误
	_ = Register(histogram) // #nosec G104
	return histogram
}

// =============================================================================
// 性能指标预设
// =============================================================================

var (
	// HTTPRequestsTotal HTTP请求总数
	HTTPRequestsTotal *Counter

	// HTTPRequestDuration HTTP请求耗时
	HTTPRequestDuration *Histogram

	// ActiveConnections 活跃连接数
	ActiveConnections *Gauge

	// MemoryUsage 内存使用
	MemoryUsage *Gauge

	// GoroutineCount Goroutine数量
	GoroutineCount *Gauge
)

func init() {
	// 初始化默认指标
	HTTPRequestsTotal = RegisterCounter(
		"http_requests_total",
		"Total number of HTTP requests",
		map[string]string{"method": "GET"},
	)

	HTTPRequestDuration = RegisterHistogram(
		"http_request_duration_seconds",
		"HTTP request latency in seconds",
		[]float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
		map[string]string{},
	)

	ActiveConnections = RegisterGauge(
		"active_connections",
		"Number of active connections",
		map[string]string{},
	)

	MemoryUsage = RegisterGauge(
		"memory_usage_bytes",
		"Memory usage in bytes",
		map[string]string{},
	)

	GoroutineCount = RegisterGauge(
		"goroutine_count",
		"Number of goroutines",
		map[string]string{},
	)
}

// UpdateRuntimeMetrics 更新运行时指标
func UpdateRuntimeMetrics() {
	// 这里可以更新Goroutine数量、内存使用等
	// 实际应用中会使用runtime包获取真实数据
	GoroutineCount.Set(int64(100))           // 示例值
	MemoryUsage.Set(int64(1024 * 1024 * 50)) // 50MB示例值
}

// StartMetricsCollector 启动指标收集器
func StartMetricsCollector(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			UpdateRuntimeMetrics()
		}
	}()
}
