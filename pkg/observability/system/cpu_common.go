package system

import "runtime"

// initLinuxCPUMonitor 初始化 Linux CPU 监控器
// 仅在 Linux 上编译，其他平台返回 nil
func initLinuxCPUMonitor() interface{} {
	if runtime.GOOS == "linux" {
		// 调用 Linux 特定的实现
		// 注意：initLinuxCPUMonitorImpl 在 cpu_linux.go 中定义（仅在 Linux 上编译）
		// 在非 Linux 平台上，cpu_other.go 提供 stub 实现
		return initLinuxCPUMonitorImpl()
	}
	return nil
}

// getLinuxCPUUsage 获取 Linux CPU 使用率
// 仅在 Linux 上编译，其他平台返回 -1
func getLinuxCPUUsage(monitor interface{}) float64 {
	if runtime.GOOS != "linux" || monitor == nil {
		return -1
	}
	// 调用 Linux 特定的实现
	// 注意：getLinuxCPUUsageImpl 在 cpu_linux.go 中定义（仅在 Linux 上编译）
	// 在非 Linux 平台上，cpu_other.go 提供 stub 实现
	return getLinuxCPUUsageImpl(monitor)
}
