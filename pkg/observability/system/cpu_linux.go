//go:build linux
// +build linux

package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CPUSnapshot CPU 快照
type CPUSnapshot struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	IOWait    uint64
	IRQ       uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
	Timestamp time.Time
}

// Total 计算总 CPU 时间
func (c *CPUSnapshot) Total() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.SoftIRQ + c.Steal + c.Guest + c.GuestNice
}

// Active 计算活跃 CPU 时间
func (c *CPUSnapshot) Active() uint64 {
	return c.User + c.Nice + c.System + c.IRQ + c.SoftIRQ + c.Steal + c.Guest + c.GuestNice
}

// linuxCPUMonitorImpl Linux 平台的精确 CPU 监控实现
type linuxCPUMonitorImpl struct {
	lastSnapshot *CPUSnapshot
}

// newLinuxCPUMonitorImpl 创建 Linux CPU 监控器
func newLinuxCPUMonitorImpl() *linuxCPUMonitorImpl {
	return &linuxCPUMonitorImpl{}
}

// initLinuxCPUMonitorImpl 初始化 Linux CPU 监控器（导出函数）
func initLinuxCPUMonitorImpl() interface{} {
	return newLinuxCPUMonitorImpl()
}

// getLinuxCPUUsageImpl 获取 Linux CPU 使用率实现（导出函数）
func getLinuxCPUUsageImpl(monitor interface{}) float64 {
	m, ok := monitor.(*linuxCPUMonitorImpl)
	if !ok {
		return -1
	}
	usage, err := m.CalculateUsage()
	if err != nil {
		return -1
	}
	return usage
}

// ReadCPUStats 读取 CPU 统计信息
func (m *linuxCPUMonitorImpl) ReadCPUStats() (*CPUSnapshot, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/stat: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read /proc/stat")
	}

	line := scanner.Text()
	if !strings.HasPrefix(line, "cpu ") {
		return nil, fmt.Errorf("invalid /proc/stat format")
	}

	fields := strings.Fields(line)
	if len(fields) < 8 {
		return nil, fmt.Errorf("insufficient fields in /proc/stat")
	}

	snapshot := &CPUSnapshot{
		Timestamp: time.Now(),
	}

	// 解析各个字段
	snapshot.User, _ = strconv.ParseUint(fields[1], 10, 64)
	snapshot.Nice, _ = strconv.ParseUint(fields[2], 10, 64)
	snapshot.System, _ = strconv.ParseUint(fields[3], 10, 64)
	snapshot.Idle, _ = strconv.ParseUint(fields[4], 10, 64)
	if len(fields) > 5 {
		snapshot.IOWait, _ = strconv.ParseUint(fields[5], 10, 64)
	}
	if len(fields) > 6 {
		snapshot.IRQ, _ = strconv.ParseUint(fields[6], 10, 64)
	}
	if len(fields) > 7 {
		snapshot.SoftIRQ, _ = strconv.ParseUint(fields[7], 10, 64)
	}
	if len(fields) > 8 {
		snapshot.Steal, _ = strconv.ParseUint(fields[8], 10, 64)
	}
	if len(fields) > 9 {
		snapshot.Guest, _ = strconv.ParseUint(fields[9], 10, 64)
	}
	if len(fields) > 10 {
		snapshot.GuestNice, _ = strconv.ParseUint(fields[10], 10, 64)
	}

	return snapshot, nil
}

// CalculateUsage 计算 CPU 使用率
func (m *linuxCPUMonitorImpl) CalculateUsage() (float64, error) {
	current, err := m.ReadCPUStats()
	if err != nil {
		return 0, err
	}

	if m.lastSnapshot == nil {
		m.lastSnapshot = current
		return 0, nil // 第一次调用，无法计算使用率
	}

	// 计算时间差
	timeDelta := current.Timestamp.Sub(m.lastSnapshot.Timestamp).Seconds()
	if timeDelta <= 0 {
		return 0, fmt.Errorf("invalid time delta")
	}

	// 计算 CPU 时间差
	totalDelta := current.Total() - m.lastSnapshot.Total()
	activeDelta := current.Active() - m.lastSnapshot.Active()

	if totalDelta == 0 {
		return 0, nil
	}

	// 计算使用率（百分比）
	usage := float64(activeDelta) / float64(totalDelta) * 100.0

	m.lastSnapshot = current

	return usage, nil
}
