package system

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// AlertLevel 告警级别
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelCritical AlertLevel = "critical"
)

// Alert 告警
type Alert struct {
	ID          string
	Level       AlertLevel
	Message     string
	MetricName  string
	MetricValue interface{}
	Threshold   interface{}
	Timestamp   time.Time
	Attributes  map[string]string
}

// AlertRule 告警规则
type AlertRule struct {
	ID          string
	Name        string
	MetricName  string
	Condition   string // "gt", "lt", "eq", "gte", "lte"
	Threshold   float64
	Level       AlertLevel
	Enabled     bool
	Duration    time.Duration // 持续时间（超过此时间才触发）
	Cooldown    time.Duration // 冷却时间（触发后多久才能再次触发）
}

// AlertHandler 告警处理器
type AlertHandler interface {
	HandleAlert(ctx context.Context, alert Alert) error
}

// AlertManager 告警管理器
type AlertManager struct {
	meter        metric.Meter
	rules        []AlertRule
	handlers     []AlertHandler
	alertHistory []Alert
	mu           sync.RWMutex
	lastAlerts   map[string]time.Time // 记录上次告警时间（用于冷却）
	enabled      bool
}

// NewAlertManager 创建告警管理器
func NewAlertManager(meter metric.Meter, enabled bool) *AlertManager {
	return &AlertManager{
		meter:      meter,
		rules:      make([]AlertRule, 0),
		handlers:   make([]AlertHandler, 0),
		alertHistory: make([]Alert, 0),
		lastAlerts: make(map[string]time.Time),
		enabled:    enabled,
	}
}

// AddRule 添加告警规则
func (am *AlertManager) AddRule(rule AlertRule) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.rules = append(am.rules, rule)
}

// RemoveRule 移除告警规则
func (am *AlertManager) RemoveRule(ruleID string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i, rule := range am.rules {
		if rule.ID == ruleID {
			am.rules = append(am.rules[:i], am.rules[i+1:]...)
			break
		}
	}
}

// AddHandler 添加告警处理器
func (am *AlertManager) AddHandler(handler AlertHandler) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.handlers = append(am.handlers, handler)
}

// Check 检查指标并触发告警
func (am *AlertManager) Check(ctx context.Context, metricName string, value float64, attrs map[string]string) error {
	if !am.enabled {
		return nil
	}

	am.mu.RLock()
	rules := make([]AlertRule, len(am.rules))
	copy(rules, am.rules)
	am.mu.RUnlock()

	for _, rule := range rules {
		if !rule.Enabled || rule.MetricName != metricName {
			continue
		}

		// 检查条件
		shouldAlert := false
		switch rule.Condition {
		case "gt":
			shouldAlert = value > rule.Threshold
		case "lt":
			shouldAlert = value < rule.Threshold
		case "gte":
			shouldAlert = value >= rule.Threshold
		case "lte":
			shouldAlert = value <= rule.Threshold
		case "eq":
			shouldAlert = value == rule.Threshold
		}

		if shouldAlert {
			// 检查冷却时间
			if am.isInCooldown(rule.ID) {
				continue
			}

			// 触发告警
			alert := Alert{
				ID:          fmt.Sprintf("%s-%d", rule.ID, time.Now().Unix()),
				Level:       rule.Level,
				Message:     fmt.Sprintf("%s: %s = %.2f (threshold: %.2f)", rule.Name, metricName, value, rule.Threshold),
				MetricName:  metricName,
				MetricValue: value,
				Threshold:   rule.Threshold,
				Timestamp:   time.Now(),
				Attributes:  attrs,
			}

			if err := am.triggerAlert(ctx, alert, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

// triggerAlert 触发告警
func (am *AlertManager) triggerAlert(ctx context.Context, alert Alert, rule AlertRule) error {
	am.mu.Lock()
	am.alertHistory = append(am.alertHistory, alert)
	if len(am.alertHistory) > 1000 {
		am.alertHistory = am.alertHistory[len(am.alertHistory)-1000:]
	}
	am.lastAlerts[rule.ID] = time.Now()
	handlers := make([]AlertHandler, len(am.handlers))
	copy(handlers, am.handlers)
	am.mu.Unlock()

	// 调用所有处理器
	for _, handler := range handlers {
		if err := handler.HandleAlert(ctx, alert); err != nil {
			// 记录错误但不中断
			fmt.Printf("Error handling alert: %v\n", err)
		}
	}

	return nil
}

// isInCooldown 检查是否在冷却期
func (am *AlertManager) isInCooldown(ruleID string) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()

	lastAlert, exists := am.lastAlerts[ruleID]
	if !exists {
		return false
	}

	// 找到对应的规则
	for _, rule := range am.rules {
		if rule.ID == ruleID {
			return time.Since(lastAlert) < rule.Cooldown
		}
	}

	return false
}

// GetAlertHistory 获取告警历史
func (am *AlertManager) GetAlertHistory(count int) []Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if count <= 0 || count > len(am.alertHistory) {
		count = len(am.alertHistory)
	}

	start := len(am.alertHistory) - count
	if start < 0 {
		start = 0
	}

	result := make([]Alert, count)
	copy(result, am.alertHistory[start:])
	return result
}

// DefaultAlertRules 返回默认告警规则
func DefaultAlertRules() []AlertRule {
	return []AlertRule{
		{
			ID:         "cpu-high",
			Name:       "CPU Usage High",
			MetricName: "system.cpu.usage",
			Condition:  "gt",
			Threshold:  80.0,
			Level:      AlertLevelWarning,
			Enabled:    true,
			Duration:   5 * time.Minute,
			Cooldown:   10 * time.Minute,
		},
		{
			ID:         "memory-high",
			Name:       "Memory Usage High",
			MetricName: "system.memory.usage",
			Condition:  "gt",
			Threshold:  90.0,
			Level:      AlertLevelCritical,
			Enabled:    true,
			Duration:   5 * time.Minute,
			Cooldown:   10 * time.Minute,
		},
		{
			ID:         "disk-high",
			Name:       "Disk Usage High",
			MetricName: "system.disk.usage",
			Condition:  "gt",
			Threshold:  85.0,
			Level:      AlertLevelWarning,
			Enabled:    true,
			Duration:   5 * time.Minute,
			Cooldown:   10 * time.Minute,
		},
	}
}
