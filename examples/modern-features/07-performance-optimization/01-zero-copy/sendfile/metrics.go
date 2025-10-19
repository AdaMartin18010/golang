package sendfile

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics 性能指标收集器
type Metrics struct {
	requestCount     int64
	bytesTransferred int64
	totalRequestTime int64 // 纳秒
	errorCount       int64
	startTime        time.Time
	mu               sync.RWMutex
}

// NewMetrics 创建新的指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

// RecordRequest 记录请求
func (m *Metrics) RecordRequest(duration time.Duration) {
	atomic.AddInt64(&m.requestCount, 1)
	atomic.AddInt64(&m.totalRequestTime, int64(duration))
}

// RecordBytesTransferred 记录传输字节数
func (m *Metrics) RecordBytesTransferred(bytes int64) {
	atomic.AddInt64(&m.bytesTransferred, bytes)
}

// RecordError 记录错误
func (m *Metrics) RecordError() {
	atomic.AddInt64(&m.errorCount, 1)
}

// GetStats 获取统计信息
func (m *Metrics) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	requestCount := atomic.LoadInt64(&m.requestCount)
	totalTime := atomic.LoadInt64(&m.totalRequestTime)
	bytesTransferred := atomic.LoadInt64(&m.bytesTransferred)
	errorCount := atomic.LoadInt64(&m.errorCount)

	var avgRequestTime time.Duration
	if requestCount > 0 {
		avgRequestTime = time.Duration(totalTime / requestCount)
	}

	uptime := time.Since(m.startTime)
	var requestsPerSecond float64
	if uptime > 0 {
		requestsPerSecond = float64(requestCount) / uptime.Seconds()
	}

	var bytesPerSecond float64
	if uptime > 0 {
		bytesPerSecond = float64(bytesTransferred) / uptime.Seconds()
	}

	var errorRate float64
	if requestCount > 0 {
		errorRate = float64(errorCount) / float64(requestCount) * 100
	}

	return Stats{
		RequestCount:       requestCount,
		BytesTransferred:   bytesTransferred,
		AverageRequestTime: avgRequestTime,
		ErrorCount:         errorCount,
		ErrorRate:          errorRate,
		RequestsPerSecond:  requestsPerSecond,
		BytesPerSecond:     bytesPerSecond,
		Uptime:             uptime,
	}
}

// Stats 统计信息
type Stats struct {
	RequestCount       int64         `json:"request_count"`
	BytesTransferred   int64         `json:"bytes_transferred"`
	AverageRequestTime time.Duration `json:"average_request_time"`
	ErrorCount         int64         `json:"error_count"`
	ErrorRate          float64       `json:"error_rate"`
	RequestsPerSecond  float64       `json:"requests_per_second"`
	BytesPerSecond     float64       `json:"bytes_per_second"`
	Uptime             time.Duration `json:"uptime"`
}

// Reset 重置指标
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.StoreInt64(&m.requestCount, 0)
	atomic.StoreInt64(&m.bytesTransferred, 0)
	atomic.StoreInt64(&m.totalRequestTime, 0)
	atomic.StoreInt64(&m.errorCount, 0)
	m.startTime = time.Now()
}
