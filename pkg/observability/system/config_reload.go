package system

import (
	"context"
	"sync"
	"time"
)

// ConfigReloader 配置热重载器
type ConfigReloader struct {
	systemMonitor *SystemMonitor
	mu            sync.RWMutex
	reloadFunc    func() (SystemConfig, error)
	checkInterval time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewConfigReloader 创建配置重载器
func NewConfigReloader(monitor *SystemMonitor, reloadFunc func() (SystemConfig, error), checkInterval time.Duration) *ConfigReloader {
	ctx, cancel := context.WithCancel(context.Background())

	return &ConfigReloader{
		systemMonitor: monitor,
		reloadFunc:    reloadFunc,
		checkInterval: checkInterval,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Start 启动配置重载器
func (cr *ConfigReloader) Start() {
	go cr.reloadLoop()
}

// Stop 停止配置重载器
func (cr *ConfigReloader) Stop() {
	if cr.cancel != nil {
		cr.cancel()
	}
}

// reloadLoop 重载循环
func (cr *ConfigReloader) reloadLoop() {
	ticker := time.NewTicker(cr.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cr.ctx.Done():
			return
		case <-ticker.C:
			cr.checkAndReload()
		}
	}
}

// checkAndReload 检查并重载配置
func (cr *ConfigReloader) checkAndReload() {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	newConfig, err := cr.reloadFunc()
	if err != nil {
		// 配置加载失败，不重载
		return
	}

	// 检查配置是否有变化
	// 这里简化实现，实际应该比较配置差异
	// 如果有变化，更新监控器配置
	_ = newConfig
}

// UpdateHealthThresholds 更新健康检查阈值
func (cr *ConfigReloader) UpdateHealthThresholds(thresholds HealthThresholds) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if cr.systemMonitor.healthChecker != nil {
		cr.systemMonitor.healthChecker.thresholds = thresholds
	}
	return nil
}

// UpdateRateLimit 更新限流配置
func (cr *ConfigReloader) UpdateRateLimit(limit int64) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if cr.systemMonitor.rateLimiter != nil {
		cr.systemMonitor.rateLimiter.UpdateLimit(limit)
	}
	return nil
}
