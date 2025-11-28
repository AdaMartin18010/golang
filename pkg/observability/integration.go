package observability

import (
	"context"
	"time"

	"github.com/yourusername/golang/pkg/observability/otlp"
	"github.com/yourusername/golang/pkg/observability/system"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Observability 可观测性集成
// 整合 OTLP、系统监控等所有可观测性功能
type Observability struct {
	otlp          *otlp.EnhancedOTLP
	systemMonitor *system.SystemMonitor
	enabled       bool
}

// Config 可观测性配置
type Config struct {
	// OTLP 配置
	ServiceName       string
	ServiceVersion    string
	OTLPEndpoint      string
	OTLPInsecure      bool
	SampleRate        float64
	MetricInterval    time.Duration
	TraceBatchTimeout time.Duration
	TraceBatchSize    int

	// 系统监控配置
	EnableSystemMonitoring bool
	SystemCollectInterval  time.Duration
	EnableDiskMonitor      bool
	EnableLoadMonitor      bool
	EnableAPMMonitor       bool
	RateLimitConfig        *system.RateLimiterConfig
	HealthThresholds       system.HealthThresholds
}

// NewObservability 创建可观测性集成
func NewObservability(cfg Config) (*Observability, error) {
	// 创建 OTLP 客户端
	otlpClient, err := otlp.NewEnhancedOTLP(otlp.Config{
		ServiceName:       cfg.ServiceName,
		ServiceVersion:    cfg.ServiceVersion,
		Endpoint:          cfg.OTLPEndpoint,
		Insecure:          cfg.OTLPInsecure,
		SampleRate:        cfg.SampleRate,
		MetricInterval:    cfg.MetricInterval,
		TraceBatchTimeout: cfg.TraceBatchTimeout,
		TraceBatchSize:    cfg.TraceBatchSize,
	})
	if err != nil {
		return nil, err
	}

	var systemMonitor *system.SystemMonitor
	if cfg.EnableSystemMonitoring {
		systemMonitor, err = system.NewSystemMonitor(system.SystemConfig{
			Meter:            otlpClient.Meter("system"),
			Tracer:           otlpClient.Tracer("system"),
			Enabled:          true,
			CollectInterval:  cfg.SystemCollectInterval,
			EnableDiskMonitor: cfg.EnableDiskMonitor,
			EnableLoadMonitor: cfg.EnableLoadMonitor,
			EnableAPMMonitor:  cfg.EnableAPMMonitor,
			RateLimitConfig:   cfg.RateLimitConfig,
			HealthThresholds:  cfg.HealthThresholds,
		})
		if err != nil {
			// 系统监控失败不影响整体功能
			systemMonitor = nil
		}
	}

	return &Observability{
		otlp:          otlpClient,
		systemMonitor: systemMonitor,
		enabled:       true,
	}, nil
}

// Start 启动所有可观测性功能
func (o *Observability) Start() error {
	if !o.enabled {
		return nil
	}

	if o.systemMonitor != nil {
		if err := o.systemMonitor.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Stop 停止所有可观测性功能
func (o *Observability) Stop(ctx context.Context) error {
	if o.systemMonitor != nil {
		if err := o.systemMonitor.Stop(); err != nil {
			return err
		}
	}

	if err := o.otlp.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Tracer 获取追踪器
func (o *Observability) Tracer(name string) trace.Tracer {
	return o.otlp.Tracer(name)
}

// Meter 获取指标器
func (o *Observability) Meter(name string) metric.Meter {
	return o.otlp.Meter(name)
}

// GetSystemMonitor 获取系统监控器
func (o *Observability) GetSystemMonitor() *system.SystemMonitor {
	return o.systemMonitor
}

// GetPlatformInfo 获取平台信息
func (o *Observability) GetPlatformInfo() system.PlatformInfo {
	if o.systemMonitor != nil {
		return o.systemMonitor.GetPlatformInfo()
	}
	return system.PlatformInfo{}
}

// IsContainer 检查是否在容器中
func (o *Observability) IsContainer() bool {
	if o.systemMonitor != nil {
		return o.systemMonitor.IsContainer()
	}
	return false
}

// IsKubernetes 检查是否在 Kubernetes 中
func (o *Observability) IsKubernetes() bool {
	if o.systemMonitor != nil {
		return o.systemMonitor.IsKubernetes()
	}
	return false
}

// IsVirtualized 检查是否在虚拟化环境中
func (o *Observability) IsVirtualized() bool {
	if o.systemMonitor != nil {
		return o.systemMonitor.IsVirtualized()
	}
	return false
}
