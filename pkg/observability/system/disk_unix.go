//go:build !windows
// +build !windows

package system

import (
	"context"
	"runtime"
	"syscall"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// collectMetricsUnix Unix 平台的磁盘监控实现
func (m *DiskMonitor) collectMetricsUnix(ctx context.Context, obs metric.Observer) error {
	// 获取当前工作目录的磁盘使用情况
	// 实际应该监控所有挂载点或特定目录
	path := "/"
	if runtime.GOOS == "darwin" {
		path = "."
	}

	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		// 如果无法获取，使用默认值
		return nil
	}

	// 计算磁盘空间
	total := int64(stat.Blocks) * int64(stat.Bsize)
	available := int64(stat.Bavail) * int64(stat.Bsize)
	used := total - available

	obs.ObserveInt64(m.diskUsageGauge, used, metric.WithAttributes(
		attribute.String("path", path),
	))
	obs.ObserveInt64(m.diskTotalGauge, total, metric.WithAttributes(
		attribute.String("path", path),
	))
	obs.ObserveInt64(m.diskAvailableGauge, available, metric.WithAttributes(
		attribute.String("path", path),
	))

	return nil
}
