package system

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DashboardData 仪表板数据
type DashboardData struct {
	Timestamp    time.Time              `json:"timestamp"`
	SystemHealth HealthStatus            `json:"system_health"`
	Metrics      map[string]interface{} `json:"metrics"`
	Alerts       []Alert                `json:"alerts"`
	Predictions  []PredictionResult     `json:"predictions"`
	PlatformInfo PlatformInfo           `json:"platform_info"`
	K8sInfo      KubernetesInfo         `json:"k8s_info,omitempty"`
}

// DashboardExporter 仪表板数据导出器
type DashboardExporter struct {
	systemMonitor *SystemMonitor
}

// NewDashboardExporter 创建仪表板导出器
func NewDashboardExporter(systemMonitor *SystemMonitor) *DashboardExporter {
	return &DashboardExporter{
		systemMonitor: systemMonitor,
	}
}

// Export 导出仪表板数据
func (de *DashboardExporter) Export(ctx context.Context) (*DashboardData, error) {
	data := &DashboardData{
		Timestamp:  time.Now(),
		Metrics:   make(map[string]interface{}),
		Alerts:     make([]Alert, 0),
		Predictions: make([]PredictionResult, 0),
	}

	// 获取健康状态
	if healthChecker := de.systemMonitor.GetHealthChecker(); healthChecker != nil {
		data.SystemHealth = healthChecker.Check(ctx)
	}

	// 获取平台信息
	data.PlatformInfo = de.systemMonitor.GetPlatformInfo()

	// 获取 Kubernetes 信息
	if k8sMonitor := de.systemMonitor.GetKubernetesMonitor(); k8sMonitor != nil && k8sMonitor.IsEnabled() {
		data.K8sInfo = k8sMonitor.GetInfo()
	}

	// 获取指标快照
	if exporter := de.systemMonitor.GetMetricsExporter(); exporter != nil {
		snapshot, err := exporter.Export(ctx)
		if err == nil {
			// 转换指标为 map
			for name, metric := range snapshot.Metrics {
				data.Metrics[name] = metric.Value
			}
		}
	}

	// 获取最近的告警
	if alertManager := de.systemMonitor.GetAlertManager(); alertManager != nil {
		data.Alerts = alertManager.GetAlertHistory(10) // 最近 10 条
	}

	// 获取预测结果
	if predictor := de.systemMonitor.GetPredictor(); predictor != nil {
		// 预测关键指标
		metrics := []string{"system.cpu.usage", "system.memory.usage", "system.disk.usage"}
		for _, metricName := range metrics {
			prediction, err := predictor.Predict(ctx, metricName, 1*time.Hour)
			if err == nil {
				data.Predictions = append(data.Predictions, *prediction)
			}
		}
	}

	return data, nil
}

// ExportJSON 导出为 JSON
func (de *DashboardExporter) ExportJSON(ctx context.Context) ([]byte, error) {
	data, err := de.Export(ctx)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(data, "", "  ")
}

// ExportForPrometheus 导出为 Prometheus 格式
func (de *DashboardExporter) ExportForPrometheus(ctx context.Context) (string, error) {
	data, err := de.Export(ctx)
	if err != nil {
		return "", err
	}

	var result string
	result += fmt.Sprintf("# HELP system_health System health status (0=unhealthy, 1=healthy)\n")
	result += fmt.Sprintf("# TYPE system_health gauge\n")
	healthValue := 0
	if data.SystemHealth.Healthy {
		healthValue = 1
	}
	result += fmt.Sprintf("system_health %d\n", healthValue)

	// 导出指标
	for name, value := range data.Metrics {
		result += fmt.Sprintf("# HELP %s Metric value\n", name)
		result += fmt.Sprintf("# TYPE %s gauge\n", name)
		result += fmt.Sprintf("%s %v\n", name, value)
	}

	return result, nil
}
