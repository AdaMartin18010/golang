package system

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// SystemMonitor 系统监控器集合
// 整合所有系统监控功能
type SystemMonitor struct {
	monitor         *Monitor
	ioMonitor       *IOMonitor
	networkMonitor  *NetworkMonitor
	diskMonitor     *DiskMonitor
	loadMonitor     *LoadMonitor
	apmMonitor      *APMMonitor
	rateLimiter      *RateLimiter
	platformMonitor  *PlatformMonitor
	kubernetesMonitor *KubernetesMonitor
	healthChecker    *HealthChecker
	aggregator       *Aggregator
	configReloader   *ConfigReloader
	metricsExporter  *MetricsExporter
	alertManager     *AlertManager
	diagnostics      *Diagnostics
	predictor        *ResourcePredictor
	enabled          bool
}

// SystemConfig 系统监控配置
type SystemConfig struct {
	Meter            metric.Meter
	Tracer           trace.Tracer // APM 需要 Tracer
	Enabled          bool
	CollectInterval  time.Duration
	EnableDiskMonitor bool // 是否启用磁盘监控
	EnableLoadMonitor  bool // 是否启用负载监控
	EnableAPMMonitor   bool // 是否启用 APM 监控
	HealthThresholds HealthThresholds // 健康检查阈值
	RateLimitConfig  *RateLimiterConfig // 限流器配置
}

// NewSystemMonitor 创建系统监控器
func NewSystemMonitor(cfg SystemConfig) (*SystemMonitor, error) {
	if cfg.Meter == nil {
		return nil, fmt.Errorf("meter is required")
	}

	// 创建平台监控器
	platformMonitor, err := NewPlatformMonitor(cfg.Meter)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform monitor: %w", err)
	}

	// 创建资源监控器
	monitor, err := NewMonitor(Config{
		Meter:           cfg.Meter,
		Enabled:         cfg.Enabled,
		CollectInterval: cfg.CollectInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create monitor: %w", err)
	}

	// 创建 IO 监控器
	ioMonitor, err := NewIOMonitor(Config{
		Meter:           cfg.Meter,
		Enabled:         cfg.Enabled,
		CollectInterval: cfg.CollectInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create IO monitor: %w", err)
	}

	// 创建网络监控器
	networkMonitor, err := NewNetworkMonitor(Config{
		Meter:           cfg.Meter,
		Enabled:         cfg.Enabled,
		CollectInterval: cfg.CollectInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create network monitor: %w", err)
	}

	// 创建磁盘监控器（可选）
	var diskMonitor *DiskMonitor
	if cfg.EnableDiskMonitor {
		diskMonitor, err = NewDiskMonitor(Config{
			Meter:           cfg.Meter,
			Enabled:         cfg.Enabled,
			CollectInterval: cfg.CollectInterval,
		})
		if err != nil {
			// 磁盘监控失败不影响整体功能
			diskMonitor = nil
		}
	}

	// 创建负载监控器（可选）
	var loadMonitor *LoadMonitor
	if cfg.EnableLoadMonitor {
		loadMonitor, err = NewLoadMonitor(Config{
			Meter:           cfg.Meter,
			Enabled:         cfg.Enabled,
			CollectInterval: cfg.CollectInterval,
		})
		if err != nil {
			loadMonitor = nil
		}
	}

	// 创建 APM 监控器（可选）
	var apmMonitor *APMMonitor
	if cfg.EnableAPMMonitor {
		apmMonitor, err = NewAPMMonitor(APMConfig{
			Meter:           cfg.Meter,
			Tracer:          cfg.Tracer,
			Enabled:         cfg.Enabled,
			CollectInterval: cfg.CollectInterval,
		})
		if err != nil {
			apmMonitor = nil
		}
	}

	// 创建限流器（可选）
	var rateLimiter *RateLimiter
	if cfg.RateLimitConfig != nil {
		rateLimiter, err = NewRateLimiter(*cfg.RateLimitConfig)
		if err != nil {
			rateLimiter = nil
		}
	}

	// 创建健康检查器
	thresholds := cfg.HealthThresholds
	if thresholds.MaxMemoryUsage == 0 {
		thresholds = DefaultHealthThresholds()
	}
	healthChecker := NewHealthChecker(monitor, thresholds)

	// 创建指标聚合器
	aggregator := NewAggregator(cfg.Meter, cfg.Enabled)

	// 创建告警管理器
	alertManager := NewAlertManager(cfg.Meter, cfg.Enabled)
	// 添加默认告警规则
	for _, rule := range DefaultAlertRules() {
		alertManager.AddRule(rule)
	}

	// 创建 Kubernetes 监控器（如果检测到 Kubernetes 环境）
	var kubernetesMonitor *KubernetesMonitor
	if platformMonitor.IsKubernetes() {
		kubernetesMonitor, err = NewKubernetesMonitor(cfg.Meter)
		if err != nil {
			// Kubernetes 监控失败不影响整体功能
			kubernetesMonitor = nil
		}
	}

	// 创建系统监控器（先创建，稍后更新引用）
	systemMonitor := &SystemMonitor{
		monitor:         monitor,
		ioMonitor:       ioMonitor,
		networkMonitor:  networkMonitor,
		diskMonitor:     diskMonitor,
		loadMonitor:     loadMonitor,
		apmMonitor:      apmMonitor,
		rateLimiter:     rateLimiter,
		platformMonitor: platformMonitor,
		healthChecker:     healthChecker,
		aggregator:        aggregator,
		kubernetesMonitor: kubernetesMonitor,
		alertManager:      alertManager,
		enabled:           cfg.Enabled,
	}

	// 现在可以安全地创建需要 systemMonitor 引用的组件
	systemMonitor.metricsExporter = NewMetricsExporter(systemMonitor, 100)
	systemMonitor.diagnostics = NewDiagnostics(systemMonitor, cfg.Meter, cfg.Enabled)
	systemMonitor.predictor = NewResourcePredictor(systemMonitor, 100)

	return systemMonitor, nil
}

// Start 启动所有监控器
func (sm *SystemMonitor) Start() error {
	if !sm.enabled {
		return nil
	}

	// 启动资源监控
	if err := sm.monitor.Start(); err != nil {
		return fmt.Errorf("failed to start monitor: %w", err)
	}

	// 启动 IO 监控
	if err := sm.ioMonitor.Start(); err != nil {
		return fmt.Errorf("failed to start IO monitor: %w", err)
	}

	// 启动网络监控
	if err := sm.networkMonitor.Start(); err != nil {
		return fmt.Errorf("failed to start network monitor: %w", err)
	}

	// 启动磁盘监控（如果启用）
	if sm.diskMonitor != nil {
		if err := sm.diskMonitor.Start(); err != nil {
			return fmt.Errorf("failed to start disk monitor: %w", err)
		}
	}

	// 启动负载监控（如果启用）
	if sm.loadMonitor != nil {
		if err := sm.loadMonitor.Start(); err != nil {
			return fmt.Errorf("failed to start load monitor: %w", err)
		}
	}

	// 启动 APM 监控（如果启用）
	if sm.apmMonitor != nil {
		if err := sm.apmMonitor.Start(); err != nil {
			return fmt.Errorf("failed to start APM monitor: %w", err)
		}
	}

	// 记录平台指标
	if err := sm.platformMonitor.RecordPlatformMetrics(context.Background()); err != nil {
		return fmt.Errorf("failed to record platform metrics: %w", err)
	}

	// 记录 Kubernetes 指标（如果启用）
	if sm.kubernetesMonitor != nil {
		if err := sm.kubernetesMonitor.RecordKubernetesMetrics(context.Background()); err != nil {
			return fmt.Errorf("failed to record Kubernetes metrics: %w", err)
		}
	}

	return nil
}

// Stop 停止所有监控器
func (sm *SystemMonitor) Stop() error {
	var errs []error

	if err := sm.monitor.Stop(); err != nil {
		errs = append(errs, err)
	}

	if err := sm.ioMonitor.Stop(); err != nil {
		errs = append(errs, err)
	}

	if err := sm.networkMonitor.Stop(); err != nil {
		errs = append(errs, err)
	}

	if sm.diskMonitor != nil {
		if err := sm.diskMonitor.Stop(); err != nil {
			errs = append(errs, err)
		}
	}

	if sm.loadMonitor != nil {
		if err := sm.loadMonitor.Stop(); err != nil {
			errs = append(errs, err)
		}
	}

	if sm.apmMonitor != nil {
		if err := sm.apmMonitor.Stop(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping monitors: %v", errs)
	}

	return nil
}

// GetPlatformInfo 获取平台信息
func (sm *SystemMonitor) GetPlatformInfo() PlatformInfo {
	return sm.platformMonitor.GetInfo()
}

// GetPlatformAttributes 获取平台属性
func (sm *SystemMonitor) GetPlatformAttributes() []interface{} {
	attrs := sm.platformMonitor.GetAttributes()
	result := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		result = append(result, attr.Key, attr.Value.AsString())
	}
	return result
}

// GetMemoryStats 获取内存统计
func (sm *SystemMonitor) GetMemoryStats() MemoryStats {
	return sm.monitor.GetMemoryStats()
}

// GetGoroutineCount 获取 Goroutine 数量
func (sm *SystemMonitor) GetGoroutineCount() int {
	return sm.monitor.GetGoroutineCount()
}

// IsContainer 检查是否在容器中
func (sm *SystemMonitor) IsContainer() bool {
	return sm.platformMonitor.IsContainer()
}

// IsKubernetes 检查是否在 Kubernetes 中
func (sm *SystemMonitor) IsKubernetes() bool {
	return sm.platformMonitor.IsKubernetes()
}

// IsVirtualized 检查是否在虚拟化环境中
func (sm *SystemMonitor) IsVirtualized() bool {
	return sm.platformMonitor.IsVirtualized()
}

// GetHealthChecker 获取健康检查器
func (sm *SystemMonitor) GetHealthChecker() *HealthChecker {
	return sm.healthChecker
}

// CheckHealth 执行健康检查
func (sm *SystemMonitor) CheckHealth(ctx context.Context) HealthStatus {
	return sm.healthChecker.Check(ctx)
}

// IsHealthy 检查是否健康
func (sm *SystemMonitor) IsHealthy(ctx context.Context) bool {
	return sm.healthChecker.IsHealthy(ctx)
}

// GetLoadMonitor 获取负载监控器
func (sm *SystemMonitor) GetLoadMonitor() *LoadMonitor {
	return sm.loadMonitor
}

// GetAPMMonitor 获取 APM 监控器
func (sm *SystemMonitor) GetAPMMonitor() *APMMonitor {
	return sm.apmMonitor
}

// GetRateLimiter 获取限流器
func (sm *SystemMonitor) GetRateLimiter() *RateLimiter {
	return sm.rateLimiter
}

// GetAggregator 获取指标聚合器
func (sm *SystemMonitor) GetAggregator() *Aggregator {
	return sm.aggregator
}

// SetConfigReloader 设置配置重载器
func (sm *SystemMonitor) SetConfigReloader(reloader *ConfigReloader) {
	sm.configReloader = reloader
}

// GetConfigReloader 获取配置重载器
func (sm *SystemMonitor) GetConfigReloader() *ConfigReloader {
	return sm.configReloader
}

// GetKubernetesMonitor 获取 Kubernetes 监控器
func (sm *SystemMonitor) GetKubernetesMonitor() *KubernetesMonitor {
	return sm.kubernetesMonitor
}

// GetKubernetesInfo 获取 Kubernetes 信息
func (sm *SystemMonitor) GetKubernetesInfo() KubernetesInfo {
	if sm.kubernetesMonitor != nil {
		return sm.kubernetesMonitor.GetInfo()
	}
	return KubernetesInfo{}
}

// GetMetricsExporter 获取指标导出器
func (sm *SystemMonitor) GetMetricsExporter() *MetricsExporter {
	return sm.metricsExporter
}

// GetAlertManager 获取告警管理器
func (sm *SystemMonitor) GetAlertManager() *AlertManager {
	return sm.alertManager
}

// GetDiagnostics 获取诊断工具
func (sm *SystemMonitor) GetDiagnostics() *Diagnostics {
	return sm.diagnostics
}

// GetPredictor 获取资源预测器
func (sm *SystemMonitor) GetPredictor() *ResourcePredictor {
	return sm.predictor
}
