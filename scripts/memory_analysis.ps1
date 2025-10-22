# Memory Analysis Tool for Go Projects
# 内存分析工具

param(
    [string]$Module = "all",
    [string]$OutputDir = "reports/memory",
    [switch]$Detailed = $false,
    [switch]$Profile = $false
)

Write-Host ""
Write-Host "╔══════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                              ║" -ForegroundColor Cyan
Write-Host "║      💾 Memory Analysis Tool v2.0  ██        ║" -ForegroundColor Yellow
Write-Host "║                                              ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# 创建输出目录
if (!(Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
}

$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$reportFile = "$OutputDir/memory-analysis-$timestamp.md"

# 函数：分析单个模块
function Analyze-Module {
    param(
        [string]$ModulePath,
        [string]$ModuleName
    )

    Write-Host "📊 Analyzing $ModuleName..." -ForegroundColor Cyan
    
    Push-Location $ModulePath
    
    try {
        # 运行带内存分析的测试
        Write-Host "  ├─ Running tests with memory profiling..." -ForegroundColor Gray
        $testOutput = go test -memprofile=mem.prof -memprofilerate=1 -run=. ./... 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ├─ ✅ Tests passed" -ForegroundColor Green
        } else {
            Write-Host "  ├─ ⚠️  Some tests failed" -ForegroundColor Yellow
        }
        
        # 分析内存profile
        if (Test-Path "mem.prof") {
            Write-Host "  ├─ Analyzing memory profile..." -ForegroundColor Gray
            
            # Top allocations
            $topAllocs = go tool pprof -top -alloc_space mem.prof 2>&1 | Select-Object -First 20
            
            # Inuse memory
            $inuseMemory = go tool pprof -top -inuse_space mem.prof 2>&1 | Select-Object -First 20
            
            Write-Host "  ├─ ✅ Profile analyzed" -ForegroundColor Green
            
            # 清理
            Remove-Item mem.prof -Force
            
            return @{
                Module = $ModuleName
                TopAllocs = $topAllocs
                InuseMemory = $inuseMemory
                Success = $true
            }
        } else {
            Write-Host "  └─ ⚠️  No profile generated" -ForegroundColor Yellow
            return @{
                Module = $ModuleName
                Success = $false
                Error = "No profile file"
            }
        }
    }
    catch {
        Write-Host "  └─ ❌ Error: $_" -ForegroundColor Red
        return @{
            Module = $ModuleName
            Success = $false
            Error = $_.Exception.Message
        }
    }
    finally {
        Pop-Location
    }
}

# 函数：运行基准测试
function Run-Benchmarks {
    param(
        [string]$ModulePath,
        [string]$ModuleName
    )

    Write-Host "⚡ Running benchmarks for $ModuleName..." -ForegroundColor Yellow
    
    Push-Location $ModulePath
    
    try {
        $benchOutput = go test -bench=. -benchmem -benchtime=5s ./... 2>&1
        
        Write-Host "  └─ ✅ Benchmarks completed" -ForegroundColor Green
        
        Pop-Location
        return @{
            Module = $ModuleName
            Output = $benchOutput
        }
    }
    catch {
        Write-Host "  └─ ❌ Error: $_" -ForegroundColor Red
        Pop-Location
        return @{
            Module = $ModuleName
            Error = $_.Exception.Message
        }
    }
}

# 函数：获取内存统计
function Get-MemoryStats {
    Write-Host "📈 Collecting memory statistics..." -ForegroundColor Cyan
    
    $modules = @(
        @{Path="pkg/agent"; Name="pkg/agent"},
        @{Path="pkg/concurrency"; Name="pkg/concurrency"},
        @{Path="pkg/http3"; Name="pkg/http3"},
        @{Path="pkg/memory"; Name="pkg/memory"},
        @{Path="pkg/observability"; Name="pkg/observability"}
    )
    
    $stats = @()
    
    foreach ($mod in $modules) {
        if (Test-Path $mod.Path) {
            Push-Location $mod.Path
            
            # 获取代码行数
            $goFiles = Get-ChildItem -Path . -Filter *.go -Recurse | Where-Object { $_.FullName -notmatch "test\.go$" }
            $lines = ($goFiles | Get-Content | Measure-Object -Line).Lines
            
            # 获取测试文件
            $testFiles = Get-ChildItem -Path . -Filter *_test.go -Recurse
            $testLines = ($testFiles | Get-Content | Measure-Object -Line).Lines
            
            Pop-Location
            
            $stats += @{
                Module = $mod.Name
                CodeLines = $lines
                TestLines = $testLines
                Files = $goFiles.Count
                TestFiles = $testFiles.Count
            }
        }
    }
    
    return $stats
}

# 开始分析
Write-Host "🔍 Starting memory analysis..." -ForegroundColor White
Write-Host ""

# 生成报告头
$report = @"
# 💾 内存分析报告

> **生成时间**: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")  
> **分析工具**: Memory Analysis Tool v2.0  
> **Go版本**: $(go version)

---

## 📋 执行摘要

本报告分析了项目各模块的内存使用情况，包括内存分配、使用模式和优化建议。

"@

# 收集统计信息
$stats = Get-MemoryStats

$report += @"

### 模块概况

| 模块 | 代码行数 | 测试行数 | 文件数 | 测试文件数 |
|------|---------|---------|--------|-----------|
"@

foreach ($stat in $stats) {
    $report += "`n| $($stat.Module) | $($stat.CodeLines) | $($stat.TestLines) | $($stat.Files) | $($stat.TestFiles) |"
}

$report += @"


---

## 🔬 详细分析

"@

# 分析各个模块
$modules = @(
    @{Path="pkg/memory"; Name="pkg/memory"},
    @{Path="pkg/http3"; Name="pkg/http3"},
    @{Path="pkg/agent"; Name="pkg/agent"},
    @{Path="pkg/concurrency"; Name="pkg/concurrency"},
    @{Path="pkg/observability"; Name="pkg/observability"}
)

$analysisResults = @()

foreach ($mod in $modules) {
    if ($Module -eq "all" -or $Module -eq $mod.Name) {
        if (Test-Path $mod.Path) {
            $result = Analyze-Module -ModulePath $mod.Path -ModuleName $mod.Name
            $analysisResults += $result
            
            if ($Detailed) {
                $benchResult = Run-Benchmarks -ModulePath $mod.Path -ModuleName $mod.Name
            }
        }
    }
}

# 添加分析结果到报告
foreach ($result in $analysisResults) {
    if ($result.Success) {
        $report += @"

### $($result.Module)

#### 内存分配 Top 10

``````text
$($result.TopAllocs | Select-Object -First 10 | Out-String)
``````

#### 使用中内存 Top 10

``````text
$($result.InuseMemory | Select-Object -First 10 | Out-String)
``````

"@
    }
}

# 添加优化建议
$report += @"

---

## 💡 优化建议

### 1. 高频分配优化

**问题**: 在热路径上频繁分配小对象

**建议**:
- 使用 `sync.Pool` 复用对象
- 预分配切片容量
- 减少字符串拼接，使用 `strings.Builder`

**示例**:
``````go
// ❌ 频繁分配
func process() {
    data := make([]byte, 1024)  // 每次都分配
    // ...
}

// ✅ 使用对象池
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    data := bufferPool.Get().([]byte)
    defer bufferPool.Put(data)
    // ...
}
``````

### 2. 减少逃逸分析

**问题**: 局部变量逃逸到堆上

**建议**:
- 使用 `go build -gcflags="-m"` 检查逃逸
- 避免在闭包中捕获大对象
- 返回值使用指针而非接口

**检查命令**:
``````bash
go build -gcflags="-m -m" ./... 2>&1 | grep "escapes to heap"
``````

### 3. 优化切片使用

**问题**: 切片容量不足导致频繁扩容

**建议**:
- 预分配已知大小的切片
- 使用 `make([]T, 0, cap)` 指定容量
- 避免切片泄漏（保留大切片的小切片）

**示例**:
``````go
// ❌ 未预分配
result := []int{}
for i := 0; i < 1000; i++ {
    result = append(result, i)  // 多次扩容
}

// ✅ 预分配
result := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    result = append(result, i)  // 无需扩容
}
``````

### 4. 避免内存泄漏

**问题**: Goroutine泄漏、闭包引用、全局变量累积

**建议**:
- 确保Goroutine能够退出
- 使用Context控制生命周期
- 定期清理缓存和映射表
- 使用弱引用模式（如有需要）

**检查工具**:
``````bash
# 运行pprof分析
go test -memprofile=mem.prof
go tool pprof -http=:8080 mem.prof

# 检查Goroutine泄漏
go test -run=TestLeak -timeout=30s
``````

### 5. 使用专用内存管理

**模块**: `pkg/memory`

**工具**:
- `GenericPool[T]` - 通用对象池
- `BytePool` - 零分配字节池
- `PoolManager` - 统一池管理
- `MemoryMonitor` - 内存监控

**示例**:
``````go
import "github.com/yourusername/golang/pkg/memory"

// 创建对象池
pool := memory.NewGenericPool(
    func() *MyObject { return &MyObject{} },
    func(obj *MyObject) { obj.Reset() },
    1000,
)

// 使用
obj := pool.Get()
defer pool.Put(obj)
// 使用obj...
``````

---

## 📊 性能基准

### 推荐的内存优化目标

| 指标 | 当前 | 目标 | 状态 |
|------|------|------|------|
| 平均分配/op | - | <1000 B/op | 待测量 |
| 分配次数/op | - | <10 allocs/op | 待测量 |
| GC暂停时间 | - | <1ms | 待测量 |
| 内存峰值 | - | <100MB | 待测量 |

### 运行基准测试

``````bash
# 所有模块
go test -bench=. -benchmem ./...

# 特定模块
go test -bench=. -benchmem ./pkg/memory/...

# 生成性能profile
go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof -http=:8080 mem.prof
``````

---

## 🔗 相关资源

- [Go内存管理官方文档](https://go.dev/doc/effective_go#allocation)
- [内存分析工具pprof](https://pkg.go.dev/net/http/pprof)
- [项目内存模块](../../pkg/memory/README.md)
- [性能优化指南](../../docs/07-性能优化/内存优化.md)

---

## 📈 下一步

1. **立即**: 修复高频分配点
2. **短期**: 实施对象池优化
3. **中期**: 建立内存监控
4. **长期**: 持续性能优化

---

**报告生成完成！** 🎉

查看详细的profile文件可以获得更多洞察。建议定期运行内存分析以保持最佳性能。

"@

# 保存报告
$report | Out-File -FilePath $reportFile -Encoding UTF8

Write-Host ""
Write-Host "✅ Analysis complete!" -ForegroundColor Green
Write-Host "📄 Report saved to: $reportFile" -ForegroundColor Cyan
Write-Host ""

# 打开报告（可选）
if ($env:TERM_PROGRAM -ne "vscode") {
    $openReport = Read-Host "Open report? (y/n)"
    if ($openReport -eq "y") {
        Start-Process $reportFile
    }
}

Write-Host "💡 Tip: Use 'go tool pprof' for interactive analysis" -ForegroundColor Yellow
Write-Host "   Example: go tool pprof -http=:8080 mem.prof" -ForegroundColor Gray
Write-Host ""

