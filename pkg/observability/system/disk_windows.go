//go:build windows
// +build windows

package system

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// collectMetricsWindows Windows 平台的磁盘监控实现
func (m *DiskMonitor) collectMetricsWindows(ctx context.Context, obs metric.Observer) error {
	// Windows 上暂时不实现磁盘监控
	// 可以使用 win32 API 或 WMI 实现
	return nil
}

// collectMetricsUnix Windows 平台的 stub（避免编译错误）
func (m *DiskMonitor) collectMetricsUnix(ctx context.Context, obs metric.Observer) error {
	// Windows 上不实现，返回 nil
	return nil
}
