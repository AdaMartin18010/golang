package observability

import (
	"time"

	"github.com/yourusername/golang/internal/config"
	"github.com/yourusername/golang/pkg/observability/system"
)

// ConfigFromAppConfig 从应用配置创建可观测性配置
func ConfigFromAppConfig(appConfig *config.Config) Config {
	obsConfig := Config{
		ServiceName:       appConfig.Observability.OTLP.ServiceName,
		ServiceVersion:    appConfig.Observability.OTLP.ServiceVersion,
		OTLPEndpoint:      appConfig.Observability.OTLP.Endpoint,
		OTLPInsecure:      appConfig.Observability.OTLP.Insecure,
		SampleRate:        0.5, // 默认值
		MetricInterval:     10 * time.Second,
		TraceBatchTimeout: 5 * time.Second,
		TraceBatchSize:    512,
	}

	// 系统监控配置
	if appConfig.Observability.System.Enabled {
		obsConfig.EnableSystemMonitoring = true
		
		// 解析收集间隔
		if appConfig.Observability.System.CollectInterval != "" {
			if interval, err := time.ParseDuration(appConfig.Observability.System.CollectInterval); err == nil {
				obsConfig.SystemCollectInterval = interval
			} else {
				obsConfig.SystemCollectInterval = 5 * time.Second
			}
		} else {
			obsConfig.SystemCollectInterval = 5 * time.Second
		}

		obsConfig.EnableDiskMonitor = appConfig.Observability.System.EnableDiskMonitor
		obsConfig.EnableLoadMonitor = appConfig.Observability.System.EnableLoadMonitor
		obsConfig.EnableAPMMonitor = appConfig.Observability.System.EnableAPMMonitor

		// 限流器配置
		if appConfig.Observability.System.RateLimit.Enabled {
			window := 1 * time.Second
			if appConfig.Observability.System.RateLimit.Window != "" {
				if w, err := time.ParseDuration(appConfig.Observability.System.RateLimit.Window); err == nil {
					window = w
				}
			}
			obsConfig.RateLimitConfig = &system.RateLimiterConfig{
				Enabled: true,
				Limit:   appConfig.Observability.System.RateLimit.Limit,
				Window:  window,
			}
		}

		// 健康检查阈值
		obsConfig.HealthThresholds = system.HealthThresholds{
			MaxMemoryUsage: appConfig.Observability.System.HealthThresholds.MaxMemoryUsage,
			MaxCPUUsage:    appConfig.Observability.System.HealthThresholds.MaxCPUUsage,
			MaxGoroutines:  appConfig.Observability.System.HealthThresholds.MaxGoroutines,
		}
	}

	return obsConfig
}

// ApplyAlertRules 应用告警规则
func ApplyAlertRules(obs *Observability, alertConfigs []config.AlertRuleConfig) {
	alertManager := obs.GetAlertManager()
	if alertManager == nil {
		return
	}

	for _, alertConfig := range alertConfigs {
		if !alertConfig.Enabled {
			continue
		}

		// 解析持续时间
		duration := 5 * time.Minute
		if alertConfig.Duration != "" {
			if d, err := time.ParseDuration(alertConfig.Duration); err == nil {
				duration = d
			}
		}

		// 解析冷却时间
		cooldown := 10 * time.Minute
		if alertConfig.Cooldown != "" {
			if c, err := time.ParseDuration(alertConfig.Cooldown); err == nil {
				cooldown = c
			}
		}

		// 转换告警级别
		var level system.AlertLevel
		switch alertConfig.Level {
		case "info":
			level = system.AlertLevelInfo
		case "warning":
			level = system.AlertLevelWarning
		case "critical":
			level = system.AlertLevelCritical
		default:
			level = system.AlertLevelWarning
		}

		// 添加告警规则
		alertManager.AddRule(system.AlertRule{
			ID:         alertConfig.ID,
			Name:       alertConfig.Name,
			MetricName: alertConfig.MetricName,
			Condition:  alertConfig.Condition,
			Threshold:  alertConfig.Threshold,
			Level:      level,
			Enabled:    alertConfig.Enabled,
			Duration:   duration,
			Cooldown:   cooldown,
		})
	}
}
