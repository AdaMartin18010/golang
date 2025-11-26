package sysinfo

import (
	"testing"
	"time"
)

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()
	if info.OS == "" {
		t.Error("Expected non-empty OS")
	}
	if info.Arch == "" {
		t.Error("Expected non-empty Arch")
	}
	if info.GoVersion == "" {
		t.Error("Expected non-empty GoVersion")
	}
	if info.NumCPU <= 0 {
		t.Error("Expected positive NumCPU")
	}
}

func TestGetMemoryInfo(t *testing.T) {
	info := GetMemoryInfo()
	if info == nil {
		t.Error("Expected non-nil MemoryInfo")
	}
}

func TestGetCPUInfo(t *testing.T) {
	info := GetCPUInfo()
	if info.NumCPU <= 0 {
		t.Error("Expected positive NumCPU")
	}
}

func TestGetOS(t *testing.T) {
	os := GetOS()
	if os == "" {
		t.Error("Expected non-empty OS")
	}
}

func TestGetArch(t *testing.T) {
	arch := GetArch()
	if arch == "" {
		t.Error("Expected non-empty Arch")
	}
}

func TestGetGoVersion(t *testing.T) {
	version := GetGoVersion()
	if version == "" {
		t.Error("Expected non-empty GoVersion")
	}
}

func TestGetNumCPU(t *testing.T) {
	num := GetNumCPU()
	if num <= 0 {
		t.Error("Expected positive NumCPU")
	}
}

func TestGetNumGoroutine(t *testing.T) {
	num := GetNumGoroutine()
	if num < 0 {
		t.Error("Expected non-negative NumGoroutine")
	}
}

func TestGetGOMAXPROCS(t *testing.T) {
	num := GetGOMAXPROCS()
	if num <= 0 {
		t.Error("Expected positive GOMAXPROCS")
	}
}

func TestSetGOMAXPROCS(t *testing.T) {
	old := GetGOMAXPROCS()
	new := SetGOMAXPROCS(old)
	if new != old {
		t.Errorf("Expected %d, got %d", old, new)
	}
}

func TestGetAllocMemory(t *testing.T) {
	mem := GetAllocMemory()
	if mem < 0 {
		t.Error("Expected non-negative memory")
	}
}

func TestFormatBytes(t *testing.T) {
	if FormatBytes(1024) == "" {
		t.Error("Expected non-empty formatted bytes")
	}
	if FormatBytes(1024*1024) == "" {
		t.Error("Expected non-empty formatted bytes")
	}
}

func TestGetMemoryUsagePercent(t *testing.T) {
	percent := GetMemoryUsagePercent()
	if percent < 0 || percent > 100 {
		t.Errorf("Expected percent between 0 and 100, got %f", percent)
	}
}

func TestGC(t *testing.T) {
	GC() // 不应该panic
}

func TestFreeOSMemory(t *testing.T) {
	FreeOSMemory() // 不应该panic
}

func TestGetStack(t *testing.T) {
	stack := GetStack()
	if len(stack) == 0 {
		t.Error("Expected non-empty stack")
	}
}

func TestIsWindows(t *testing.T) {
	_ = IsWindows() // 不应该panic
}

func TestIsLinux(t *testing.T) {
	_ = IsLinux() // 不应该panic
}

func TestIsDarwin(t *testing.T) {
	_ = IsDarwin() // 不应该panic
}

func TestIsUnix(t *testing.T) {
	_ = IsUnix() // 不应该panic
}

func TestMonitor(t *testing.T) {
	monitor := NewMonitor(100*time.Millisecond, func(sysInfo *SystemInfo, memInfo *MemoryInfo, cpuInfo *CPUInfo) {
		// 回调函数
	})

	// 启动监控
	go monitor.Start()

	// 等待一段时间
	time.Sleep(200 * time.Millisecond)

	// 停止监控
	monitor.Stop()
}
