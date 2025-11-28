//go:build !linux
// +build !linux

package system

// initLinuxCPUMonitorImpl 非 Linux 平台的 stub
func initLinuxCPUMonitorImpl() interface{} {
	return nil
}

// getLinuxCPUUsageImpl 非 Linux 平台的 stub
func getLinuxCPUUsageImpl(monitor interface{}) float64 {
	return -1
}
