# PowerShell版本的代码质量扫描器
# Code Quality and Runability Scanner for Windows

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Go Project Quality Scanner" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# 统计变量
$totalGoFiles = 0
$totalTestFiles = 0
$compilableModules = 0
$nonCompilableModules = 0
$totalModules = 0

# 创建报告文件
$reportFile = "code_quality_report_$(Get-Date -Format 'yyyyMMdd_HHmmss').md"

"# Code Quality Scan Report" | Out-File -FilePath $reportFile
"" | Out-File -FilePath $reportFile -Append
"**Date**: $(Get-Date)" | Out-File -FilePath $reportFile -Append
"**Go Version**: $(go version)" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append

# 1. 统计Go文件
Write-Host "[1/6] Counting Go files..." -ForegroundColor Yellow
$goFiles = Get-ChildItem -Path . -Recurse -Filter "*.go" | Where-Object { $_.FullName -notmatch "vendor|\.git" }
$testFiles = $goFiles | Where-Object { $_.Name -match "_test\.go$" }
$totalGoFiles = $goFiles.Count
$totalTestFiles = $testFiles.Count
$totalCodeFiles = $totalGoFiles - $totalTestFiles

Write-Host "  - Total Go files: $totalGoFiles"
Write-Host "  - Code files: $totalCodeFiles"
Write-Host "  - Test files: $totalTestFiles"

# 2. 检查模块编译
Write-Host ""
Write-Host "[2/6] Checking module compilation..." -ForegroundColor Yellow
"" | Out-File -FilePath $reportFile -Append
"## Compilation Status" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append

$failedModules = @()
$modFiles = Get-ChildItem -Path . -Recurse -Filter "go.mod" | Where-Object { $_.FullName -notmatch "vendor" }

foreach ($modFile in $modFiles) {
    $modDir = $modFile.DirectoryName
    $totalModules++
    
    Write-Host "  Checking $modDir... " -NoNewline
    
    Push-Location $modDir
    $buildResult = go build ./... 2>&1
    $buildSuccess = $LASTEXITCODE -eq 0
    Pop-Location
    
    if ($buildSuccess) {
        Write-Host "✓" -ForegroundColor Green
        $compilableModules++
    } else {
        Write-Host "✗" -ForegroundColor Red
        $nonCompilableModules++
        $failedModules += $modDir
        "- ❌ $modDir" | Out-File -FilePath $reportFile -Append
    }
}

if ($failedModules.Count -eq 0) {
    "" | Out-File -FilePath $reportFile -Append
    "✅ **All modules compile successfully!**" | Out-File -FilePath $reportFile -Append
}

# 3. 运行测试
Write-Host ""
Write-Host "[3/6] Running tests..." -ForegroundColor Yellow
"" | Out-File -FilePath $reportFile -Append
"## Test Results" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append

$testOutput = go test ./... -cover -coverprofile=coverage.out 2>&1
$testSuccess = $LASTEXITCODE -eq 0

if ($testSuccess) {
    Write-Host "Tests passed!" -ForegroundColor Green
    
    if (Test-Path coverage.out) {
        $coverageInfo = go tool cover -func=coverage.out | Select-String "total"
        if ($coverageInfo) {
            $coverage = ($coverageInfo -split '\s+')[-1]
            "- **Coverage**: $coverage" | Out-File -FilePath $reportFile -Append
            Write-Host "  Coverage: $coverage"
        }
    }
} else {
    Write-Host "Some tests failed!" -ForegroundColor Red
    "- ⚠️ **Some tests failed**" | Out-File -FilePath $reportFile -Append
}

# 4. 代码质量检查
Write-Host ""
Write-Host "[4/6] Running code quality checks..." -ForegroundColor Yellow
"" | Out-File -FilePath $reportFile -Append
"## Code Quality" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append

# go vet
Write-Host "  Running go vet... " -NoNewline
$vetOutput = go vet ./... 2>&1
$vetSuccess = $LASTEXITCODE -eq 0

if ($vetSuccess) {
    Write-Host "✓" -ForegroundColor Green
    "- ✅ **go vet**: Passed" | Out-File -FilePath $reportFile -Append
} else {
    Write-Host "✗" -ForegroundColor Red
    "- ⚠️ **go vet**: Issues found" | Out-File -FilePath $reportFile -Append
}

# gofmt
Write-Host "  Checking formatting... " -NoNewline
$unformatted = gofmt -l -s . 2>&1 | Where-Object { $_ -match '\.go$' -and $_ -notmatch 'vendor' }

if ($null -eq $unformatted -or $unformatted.Count -eq 0) {
    Write-Host "✓" -ForegroundColor Green
    "- ✅ **gofmt**: All files properly formatted" | Out-File -FilePath $reportFile -Append
} else {
    Write-Host "✗" -ForegroundColor Red
    "- ⚠️ **gofmt**: Some files need formatting" | Out-File -FilePath $reportFile -Append
}

# 5. 依赖检查
Write-Host ""
Write-Host "[5/6] Checking dependencies..." -ForegroundColor Yellow
"" | Out-File -FilePath $reportFile -Append
"## Dependencies" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append

$moduleCount = $modFiles.Count
"- **Total modules**: $moduleCount" | Out-File -FilePath $reportFile -Append

# 6. 生成总结
Write-Host ""
Write-Host "[6/6] Generating summary..." -ForegroundColor Yellow
"" | Out-File -FilePath $reportFile -Append
"## Summary" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append
"| Metric | Value |" | Out-File -FilePath $reportFile -Append
"|--------|-------|" | Out-File -FilePath $reportFile -Append
"| Total Go files | $totalGoFiles |" | Out-File -FilePath $reportFile -Append
"| Code files | $totalCodeFiles |" | Out-File -FilePath $reportFile -Append
"| Test files | $totalTestFiles |" | Out-File -FilePath $reportFile -Append
"| Total modules | $totalModules |" | Out-File -FilePath $reportFile -Append
"| Compilable modules | $compilableModules |" | Out-File -FilePath $reportFile -Append
"| Non-compilable modules | $nonCompilableModules |" | Out-File -FilePath $reportFile -Append

$testCoveragePercent = 0
if ($totalTestFiles -gt 0 -and $totalCodeFiles -gt 0) {
    $testCoveragePercent = [math]::Round(($totalTestFiles * 100 / $totalCodeFiles), 2)
}

"| Test file coverage | $testCoveragePercent% |" | Out-File -FilePath $reportFile -Append

# 计算综合评分
$compilationScore = 0
if ($totalModules -gt 0) {
    $compilationScore = [math]::Round(($compilableModules * 100 / $totalModules), 2)
}

"" | Out-File -FilePath $reportFile -Append
"## Quality Score" | Out-File -FilePath $reportFile -Append
"" | Out-File -FilePath $reportFile -Append
"- **Compilation Success Rate**: $compilationScore%" | Out-File -FilePath $reportFile -Append
"- **Test File Ratio**: $testCoveragePercent%" | Out-File -FilePath $reportFile -Append

# 显示结果
Write-Host ""
Write-Host "===================================" -ForegroundColor Green
Write-Host "Scan Complete!" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Green
Write-Host ""
Write-Host "Summary:"
Write-Host "  - Total Go files: $totalGoFiles"
Write-Host "  - Test files: $totalTestFiles ($testCoveragePercent%)"
Write-Host "  - Compilable modules: $compilableModules/$totalModules ($compilationScore%)"
Write-Host ""
Write-Host "Report saved to: $reportFile" -ForegroundColor Cyan
Write-Host ""

# 如果有失败的模块，列出来
if ($failedModules.Count -gt 0) {
    Write-Host "Failed modules:" -ForegroundColor Red
    foreach ($mod in $failedModules) {
        Write-Host "  - $mod"
    }
    Write-Host ""
    exit 1
}

exit 0

