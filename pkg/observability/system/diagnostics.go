package system

import (
	"context"
	"encoding/json"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// DiagnosticReport 诊断报告
type DiagnosticReport struct {
	Timestamp    time.Time              `json:"timestamp"`
	SystemInfo   SystemInfo             `json:"system_info"`
	Metrics      map[string]interface{} `json:"metrics"`
	HealthStatus HealthStatus           `json:"health_status"`
	Issues       []DiagnosticIssue      `json:"issues"`
	Recommendations []string            `json:"recommendations"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	GoVersion       string `json:"go_version"`
	NumCPU          int    `json:"num_cpu"`
	NumGoroutines   int    `json:"num_goroutines"`
	MemoryAlloc     uint64 `json:"memory_alloc"`
	MemoryTotal     uint64 `json:"memory_total"`
	MemorySys       uint64 `json:"memory_sys"`
	GC              uint32 `json:"gc"`
}

// DiagnosticIssue 诊断问题
type DiagnosticIssue struct {
	Level       string `json:"level"` // info, warning, error
	Category    string `json:"category"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	Value       interface{} `json:"value"`
	Threshold   interface{} `json:"threshold"`
}

// Diagnostics 诊断工具
type Diagnostics struct {
	systemMonitor *SystemMonitor
	meter         metric.Meter
	enabled       bool
}

// NewDiagnostics 创建诊断工具
func NewDiagnostics(systemMonitor *SystemMonitor, meter metric.Meter, enabled bool) *Diagnostics {
	return &Diagnostics{
		systemMonitor: systemMonitor,
		meter:         meter,
		enabled:       enabled,
	}
}

// GenerateReport 生成诊断报告
func (d *Diagnostics) GenerateReport(ctx context.Context) (*DiagnosticReport, error) {
	report := &DiagnosticReport{
		Timestamp:  time.Now(),
		Metrics:    make(map[string]interface{}),
		Issues:     make([]DiagnosticIssue, 0),
		Recommendations: make([]string, 0),
	}

	// 收集系统信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	report.SystemInfo = SystemInfo{
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		GoVersion:     runtime.Version(),
		NumCPU:        runtime.NumCPU(),
		NumGoroutines: runtime.NumGoroutine(),
		MemoryAlloc:   m.Alloc,
		MemoryTotal:   m.TotalAlloc,
		MemorySys:     m.Sys,
		GC:            m.NumGC,
	}

	// 收集健康状态
	if healthChecker := d.systemMonitor.GetHealthChecker(); healthChecker != nil {
		status := healthChecker.Check(ctx)
		report.HealthStatus = status
	}

	// 从聚合器获取指标
	if aggregator := d.systemMonitor.GetAggregator(); aggregator != nil {
		// 获取关键指标
		// 这里简化实现
	}

	// 诊断问题
	d.diagnoseIssues(report)

	// 生成建议
	d.generateRecommendations(report)

	return report, nil
}

// diagnoseIssues 诊断问题
func (d *Diagnostics) diagnoseIssues(report *DiagnosticReport) {
	// 检查内存使用
	if report.SystemInfo.MemoryAlloc > 100*1024*1024 { // 100MB
		report.Issues = append(report.Issues, DiagnosticIssue{
			Level:       "warning",
			Category:    "memory",
			Description: "High memory allocation detected",
			Metric:      "memory.alloc",
			Value:       report.SystemInfo.MemoryAlloc,
		})
	}

	// 检查 Goroutine 数量
	if report.SystemInfo.NumGoroutines > 1000 {
		report.Issues = append(report.Issues, DiagnosticIssue{
			Level:       "warning",
			Category:    "goroutines",
			Description: "High number of goroutines detected",
			Metric:      "goroutines.count",
			Value:       report.SystemInfo.NumGoroutines,
		})
	}

	// 检查 GC 频率
	if report.SystemInfo.GC > 100 {
		report.Issues = append(report.Issues, DiagnosticIssue{
			Level:       "info",
			Category:    "gc",
			Description: "High GC frequency detected",
			Metric:      "gc.count",
			Value:       report.SystemInfo.GC,
		})
	}

	// 检查健康状态
	if !report.HealthStatus.Healthy {
		report.Issues = append(report.Issues, DiagnosticIssue{
			Level:       "error",
			Category:    "health",
			Description: "System health check failed",
			Metric:      "health.status",
			Value:       report.HealthStatus,
		})
	}
}

// generateRecommendations 生成建议
func (d *Diagnostics) generateRecommendations(report *DiagnosticReport) {
	for _, issue := range report.Issues {
		switch issue.Category {
		case "memory":
			report.Recommendations = append(report.Recommendations,
				"Consider optimizing memory usage or increasing available memory")
		case "goroutines":
			report.Recommendations = append(report.Recommendations,
				"Review goroutine usage and consider using worker pools")
		case "gc":
			report.Recommendations = append(report.Recommendations,
				"Consider optimizing object allocation patterns")
		case "health":
			report.Recommendations = append(report.Recommendations,
				"Review system health metrics and address underlying issues")
		}
	}
}

// ExportJSON 导出为 JSON
func (d *Diagnostics) ExportJSON(ctx context.Context) ([]byte, error) {
	report, err := d.GenerateReport(ctx)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(report, "", "  ")
}
