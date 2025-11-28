package system

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PredictionResult 预测结果
type PredictionResult struct {
	MetricName    string    `json:"metric_name"`
	CurrentValue  float64   `json:"current_value"`
	PredictedValue float64  `json:"predicted_value"`
	TimeHorizon   time.Duration `json:"time_horizon"`
	Confidence    float64   `json:"confidence"` // 0-1
	Trend         string    `json:"trend"`      // "increasing", "decreasing", "stable"
}

// ResourcePredictor 资源使用预测器
type ResourcePredictor struct {
	systemMonitor *SystemMonitor
	history       []MetricSnapshot
	maxHistory    int
	mu            sync.RWMutex
}

// NewResourcePredictor 创建资源预测器
func NewResourcePredictor(systemMonitor *SystemMonitor, maxHistory int) *ResourcePredictor {
	if maxHistory <= 0 {
		maxHistory = 100
	}

	return &ResourcePredictor{
		systemMonitor: systemMonitor,
		history:       make([]MetricSnapshot, 0, maxHistory),
		maxHistory:    maxHistory,
	}
}

// Predict 预测资源使用
func (rp *ResourcePredictor) Predict(ctx context.Context, metricName string, timeHorizon time.Duration) (*PredictionResult, error) {
	rp.mu.RLock()
	defer rp.mu.RUnlock()

	if len(rp.history) < 2 {
		return nil, fmt.Errorf("insufficient history for prediction")
	}

	// 获取最近的指标值
	recent := rp.getRecentValues(metricName, 10)
	if len(recent) < 2 {
		return nil, fmt.Errorf("insufficient data for metric: %s", metricName)
	}

	// 简单的线性预测
	currentValue := recent[len(recent)-1]
	previousValue := recent[len(recent)-2]

	// 计算趋势
	trend := "stable"
	slope := (currentValue - previousValue) / float64(len(recent))
	if slope > 0.1 {
		trend = "increasing"
	} else if slope < -0.1 {
		trend = "decreasing"
	}

	// 预测未来值
	predictedValue := currentValue + slope*float64(timeHorizon.Seconds())

	// 计算置信度（基于历史数据的稳定性）
	confidence := rp.calculateConfidence(recent)

	return &PredictionResult{
		MetricName:     metricName,
		CurrentValue:   currentValue,
		PredictedValue: predictedValue,
		TimeHorizon:    timeHorizon,
		Confidence:     confidence,
		Trend:          trend,
	}, nil
}

// getRecentValues 获取最近的指标值
func (rp *ResourcePredictor) getRecentValues(metricName string, count int) []float64 {
	var values []float64

	for i := len(rp.history) - 1; i >= 0 && len(values) < count; i-- {
		snapshot := rp.history[i]
		if metric, exists := snapshot.Metrics[metricName]; exists {
			if val, ok := metric.Value.(float64); ok {
				values = append([]float64{val}, values...)
			}
		}
	}

	return values
}

// calculateConfidence 计算置信度
func (rp *ResourcePredictor) calculateConfidence(values []float64) float64 {
	if len(values) < 2 {
		return 0.5
	}

	// 计算方差
	var sum, mean float64
	for _, v := range values {
		sum += v
	}
	mean = sum / float64(len(values))

	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(values))

	// 方差越小，置信度越高
	// 归一化到 0-1
	confidence := 1.0 / (1.0 + variance)
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// AddSnapshot 添加快照
func (rp *ResourcePredictor) AddSnapshot(snapshot MetricSnapshot) {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	rp.history = append(rp.history, snapshot)

	if len(rp.history) > rp.maxHistory {
		rp.history = rp.history[len(rp.history)-rp.maxHistory:]
	}
}
