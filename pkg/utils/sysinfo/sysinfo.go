package sysinfo

import (
	"runtime"
	"strconv"
	"time"
)

// SystemInfo 系统信息
type SystemInfo struct {
	OS           string    `json:"os"`
	Arch         string    `json:"arch"`
	GoVersion    string    `json:"go_version"`
	NumCPU       int       `json:"num_cpu"`
	NumGoroutine int       `json:"num_goroutine"`
	Timestamp    time.Time `json:"timestamp"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Alloc      uint64 `json:"alloc"`       // 已分配的内存
	TotalAlloc uint64 `json:"total_alloc"` // 累计分配的内存
	Sys        uint64 `json:"sys"`         // 系统内存
	NumGC      uint32 `json:"num_gc"`      // GC次数
}

// CPUInfo CPU信息
type CPUInfo struct {
	NumCPU       int     `json:"num_cpu"`        // CPU核心数
	NumGoroutine int     `json:"num_goroutine"`  // Goroutine数量
	GOMAXPROCS   int     `json:"gomaxprocs"`     // GOMAXPROCS设置
	CPUUsage     float64 `json:"cpu_usage"`      // CPU使用率（需要计算）
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		Timestamp:    time.Now(),
	}
}

// GetMemoryInfo 获取内存信息
func GetMemoryInfo() *MemoryInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &MemoryInfo{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}
}

// GetCPUInfo 获取CPU信息
func GetCPUInfo() *CPUInfo {
	return &CPUInfo{
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		GOMAXPROCS:   runtime.GOMAXPROCS(0),
	}
}

// GetOS 获取操作系统
func GetOS() string {
	return runtime.GOOS
}

// GetArch 获取架构
func GetArch() string {
	return runtime.GOARCH
}

// GetGoVersion 获取Go版本
func GetGoVersion() string {
	return runtime.Version()
}

// GetNumCPU 获取CPU核心数
func GetNumCPU() int {
	return runtime.NumCPU()
}

// GetNumGoroutine 获取Goroutine数量
func GetNumGoroutine() int {
	return runtime.NumGoroutine()
}

// GetGOMAXPROCS 获取GOMAXPROCS设置
func GetGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}

// SetGOMAXPROCS 设置GOMAXPROCS
func SetGOMAXPROCS(n int) int {
	return runtime.GOMAXPROCS(n)
}

// GetAllocMemory 获取已分配的内存（字节）
func GetAllocMemory() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// GetTotalAllocMemory 获取累计分配的内存（字节）
func GetTotalAllocMemory() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc
}

// GetSysMemory 获取系统内存（字节）
func GetSysMemory() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Sys
}

// GetNumGC 获取GC次数
func GetNumGC() uint32 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.NumGC
}

// FormatBytes 格式化字节数
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return formatFloat(float64(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return formatFloat(float64(bytes)/float64(div)) + " " + string("KMGTPE"[exp]) + "B"
}

// formatFloat 格式化浮点数
func formatFloat(f float64) string {
	if f < 10 {
		return formatFloatPrecision(f, 2)
	} else if f < 100 {
		return formatFloatPrecision(f, 1)
	}
	return formatFloatPrecision(f, 0)
}

// formatFloatPrecision 格式化浮点数到指定精度
func formatFloatPrecision(f float64, precision int) string {
	if precision == 0 {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	return strconv.FormatFloat(f, 'f', precision, 64)
}

// GetMemoryUsagePercent 获取内存使用率（百分比）
func GetMemoryUsagePercent() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if m.Sys == 0 {
		return 0
	}
	return float64(m.Alloc) / float64(m.Sys) * 100
}

// GC 执行GC
func GC() {
	runtime.GC()
}

// FreeOSMemory 释放OS内存
func FreeOSMemory() {
	runtime.GC()
	runtime.GC()
}

// GetStack 获取当前goroutine的堆栈信息
func GetStack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// GetAllStacks 获取所有goroutine的堆栈信息
func GetAllStacks() []byte {
	buf := make([]byte, 1024*1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// GetCaller 获取调用者信息
func GetCaller(skip int) (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(skip + 1)
}

// GetCallers 获取调用栈
func GetCallers(skip int, pc []uintptr) int {
	return runtime.Callers(skip+1, pc)
}

// GetFuncName 获取函数名
func GetFuncName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	return fn.Name()
}

// GetFileLine 获取文件和行号
func GetFileLine(pc uintptr) (file string, line int) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "", 0
	}
	return fn.FileLine(pc)
}

// IsWindows 检查是否为Windows系统
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux 检查是否为Linux系统
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsDarwin 检查是否为Darwin系统（macOS）
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsUnix 检查是否为Unix系统
func IsUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd"
}

// IsAMD64 检查是否为AMD64架构
func IsAMD64() bool {
	return runtime.GOARCH == "amd64"
}

// IsARM64 检查是否为ARM64架构
func IsARM64() bool {
	return runtime.GOARCH == "arm64"
}

// Is386 检查是否为386架构
func Is386() bool {
	return runtime.GOARCH == "386"
}

// GetCompiler 获取编译器信息
func GetCompiler() string {
	return runtime.Compiler
}

// GetNumCgoCall 获取CGO调用次数
func GetNumCgoCall() int64 {
	return runtime.NumCgoCall()
}

// Monitor 系统监控器
type Monitor struct {
	interval time.Duration
	stop     chan struct{}
	callback func(*SystemInfo, *MemoryInfo, *CPUInfo)
}

// NewMonitor 创建系统监控器
func NewMonitor(interval time.Duration, callback func(*SystemInfo, *MemoryInfo, *CPUInfo)) *Monitor {
	return &Monitor{
		interval: interval,
		stop:     make(chan struct{}),
		callback: callback,
	}
}

// Start 启动监控
func (m *Monitor) Start() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sysInfo := GetSystemInfo()
			memInfo := GetMemoryInfo()
			cpuInfo := GetCPUInfo()
			m.callback(sysInfo, memInfo, cpuInfo)
		case <-m.stop:
			return
		}
	}
}

// Stop 停止监控
func (m *Monitor) Stop() {
	close(m.stop)
}
