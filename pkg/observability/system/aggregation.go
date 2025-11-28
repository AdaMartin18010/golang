package system

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Aggregator 指标聚合器
// 提供指标聚合和统计功能
type Aggregator struct {
	meter    metric.Meter
	enabled  bool
	mu       sync.RWMutex
	counters map[string]int64
	gauges   map[string]float64
	histograms map[string][]float64
}

// NewAggregator 创建指标聚合器
func NewAggregator(meter metric.Meter, enabled bool) *Aggregator {
	return &Aggregator{
		meter:      meter,
		enabled:    enabled,
		counters:   make(map[string]int64),
		gauges:     make(map[string]float64),
		histograms: make(map[string][]float64),
	}
}

// RecordCounter 记录计数器值
func (a *Aggregator) RecordCounter(name string, value int64, attrs ...attribute.KeyValue) {
	if !a.enabled {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	key := a.makeKey(name, attrs)
	a.counters[key] += value
}

// RecordGauge 记录 Gauge 值
func (a *Aggregator) RecordGauge(name string, value float64, attrs ...attribute.KeyValue) {
	if !a.enabled {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	key := a.makeKey(name, attrs)
	a.gauges[key] = value
}

// RecordHistogram 记录直方图值
func (a *Aggregator) RecordHistogram(name string, value float64, attrs ...attribute.KeyValue) {
	if !a.enabled {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	key := a.makeKey(name, attrs)
	a.histograms[key] = append(a.histograms[key], value)
}

// GetCounter 获取计数器值
func (a *Aggregator) GetCounter(name string, attrs ...attribute.KeyValue) int64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	key := a.makeKey(name, attrs)
	return a.counters[key]
}

// GetGauge 获取 Gauge 值
func (a *Aggregator) GetGauge(name string, attrs ...attribute.KeyValue) float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	key := a.makeKey(name, attrs)
	return a.gauges[key]
}

// GetHistogramStats 获取直方图统计
func (a *Aggregator) GetHistogramStats(name string, attrs ...attribute.KeyValue) HistogramStats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	key := a.makeKey(name, attrs)
	values := a.histograms[key]

	if len(values) == 0 {
		return HistogramStats{}
	}

	// 计算统计信息
	var sum float64
	min := values[0]
	max := values[0]

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return HistogramStats{
		Count: len(values),
		Sum:   sum,
		Min:   min,
		Max:   max,
		Mean:  sum / float64(len(values)),
	}
}

// Reset 重置所有聚合数据
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.counters = make(map[string]int64)
	a.gauges = make(map[string]float64)
	a.histograms = make(map[string][]float64)
}

// makeKey 生成键
func (a *Aggregator) makeKey(name string, attrs []attribute.KeyValue) string {
	key := name
	for _, attr := range attrs {
		key += fmt.Sprintf(":%s=%v", attr.Key, attr.Value.AsString())
	}
	return key
}

// HistogramStats 直方图统计
type HistogramStats struct {
	Count int
	Sum   float64
	Min   float64
	Max   float64
	Mean  float64
}
