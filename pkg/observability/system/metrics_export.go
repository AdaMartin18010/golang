package system

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// MetricSnapshot 指标快照
type MetricSnapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]MetricValue `json:"metrics"`
}

// MetricValue 指标值
type MetricValue struct {
	Type      string                 `json:"type"`      // counter, gauge, histogram
	Value     interface{}            `json:"value"`      // 实际值
	Unit      string                 `json:"unit"`      // 单位
	Attributes map[string]string     `json:"attributes"` // 属性
}

// MetricsExporter 指标导出器
// 提供指标查询和导出功能
type MetricsExporter struct {
	systemMonitor *SystemMonitor
	history       []MetricSnapshot
	maxHistory    int
	mu            sync.RWMutex
}

// NewMetricsExporter 创建指标导出器
func NewMetricsExporter(systemMonitor *SystemMonitor, maxHistory int) *MetricsExporter {
	if maxHistory <= 0 {
		maxHistory = 100 // 默认保留 100 个快照
	}

	return &MetricsExporter{
		systemMonitor: systemMonitor,
		history:       make([]MetricSnapshot, 0, maxHistory),
		maxHistory:    maxHistory,
	}
}

// Export 导出当前指标快照
func (me *MetricsExporter) Export(ctx context.Context) (*MetricSnapshot, error) {
	snapshot := &MetricSnapshot{
		Timestamp: time.Now(),
		Metrics:   make(map[string]MetricValue),
	}

	// 从聚合器获取指标
	if aggregator := me.systemMonitor.GetAggregator(); aggregator != nil {
		// 导出计数器
		// 这里简化实现，实际应该从 MeterProvider 读取
	}

	return snapshot, nil
}

// ExportJSON 导出为 JSON
func (me *MetricsExporter) ExportJSON(ctx context.Context) ([]byte, error) {
	snapshot, err := me.Export(ctx)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

// GetHistory 获取历史快照
func (me *MetricsExporter) GetHistory(count int) []MetricSnapshot {
	me.mu.RLock()
	defer me.mu.RUnlock()

	if count <= 0 || count > len(me.history) {
		count = len(me.history)
	}

	start := len(me.history) - count
	if start < 0 {
		start = 0
	}

	result := make([]MetricSnapshot, count)
	copy(result, me.history[start:])
	return result
}

// SaveSnapshot 保存快照到历史
func (me *MetricsExporter) SaveSnapshot(snapshot *MetricSnapshot) {
	me.mu.Lock()
	defer me.mu.Unlock()

	me.history = append(me.history, *snapshot)

	// 限制历史记录数量
	if len(me.history) > me.maxHistory {
		me.history = me.history[len(me.history)-me.maxHistory:]
	}
}

// ClearHistory 清空历史
func (me *MetricsExporter) ClearHistory() {
	me.mu.Lock()
	defer me.mu.Unlock()
	me.history = make([]MetricSnapshot, 0, me.maxHistory)
}

// QueryMetrics 查询指标
type QueryOptions struct {
	MetricNames []string
	StartTime   time.Time
	EndTime     time.Time
	Attributes  map[string]string
}

// Query 查询指标
func (me *MetricsExporter) Query(ctx context.Context, opts QueryOptions) ([]MetricSnapshot, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	var results []MetricSnapshot

	for _, snapshot := range me.history {
		// 时间过滤
		if !opts.StartTime.IsZero() && snapshot.Timestamp.Before(opts.StartTime) {
			continue
		}
		if !opts.EndTime.IsZero() && snapshot.Timestamp.After(opts.EndTime) {
			continue
		}

		// 指标名称过滤
		if len(opts.MetricNames) > 0 {
			matched := false
			for _, name := range opts.MetricNames {
				if _, exists := snapshot.Metrics[name]; exists {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		results = append(results, snapshot)
	}

	return results, nil
}
