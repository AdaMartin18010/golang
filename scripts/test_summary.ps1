# 测试统计脚本
Write-Host "=== Go项目测试统计报告 ===" -ForegroundColor Cyan
Write-Host ""

$modules = @(
    "examples/concurrency",
    "docs/02-Go语言现代化/14-Go-1.25并发和网络/examples/waitgroup_go",
    "docs/02-Go语言现代化/08-智能化架构集成/01-AI-Agent架构"
)

$totalPass = 0
$totalFail = 0
$totalTests = 0
$startDir = Get-Location

foreach ($module in $modules) {
    Write-Host "测试模块: $module" -ForegroundColor Yellow
    
    if (Test-Path $module) {
        Set-Location $module
        $output = go test -v ./... 2>&1 | Out-String
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ PASS" -ForegroundColor Green
            $totalPass++
            
            # 统计测试数量
            $passCount = ([regex]::Matches($output, "--- PASS:")).Count
            $totalTests += $passCount
            Write-Host "  通过测试: $passCount 个" -ForegroundColor Gray
        } else {
            Write-Host "  ❌ FAIL" -ForegroundColor Red
            $totalFail++
        }
        
        Set-Location $startDir
    } else {
        Write-Host "  ⚠ 模块不存在" -ForegroundColor Yellow
    }
    
    Write-Host ""
}

Write-Host "=== 总结 ===" -ForegroundColor Cyan
Write-Host "模块通过: $totalPass / $($modules.Count)" -ForegroundColor Green
Write-Host "模块失败: $totalFail / $($modules.Count)" $(if ($totalFail -eq 0) {"-ForegroundColor Green"} else {"-ForegroundColor Red"})
Write-Host "测试用例总数: $totalTests 个" -ForegroundColor Cyan
Write-Host ""

if ($totalFail -eq 0) {
    Write-Host "🎉 所有模块测试通过！" -ForegroundColor Green
} else {
    Write-Host "⚠ 存在失败的模块" -ForegroundColor Yellow
}

