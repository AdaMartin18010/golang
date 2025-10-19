# 🎨 代码格式化脚本
# 格式化所有Go代码，运行vet检查

Write-Host "🎨 开始代码格式化..." -ForegroundColor Cyan
Write-Host ""

# 统计
$totalFiles = 0
$formattedFiles = 0
$vetIssues = 0

# 排除的目录
$excludeDirs = @(
    ".git",
    "node_modules",
    ".cursor",
    "vendor",
    "docs\00-备份"
)

# Step 1: gofmt 格式化
Write-Host "📝 Step 1: 运行 gofmt..." -ForegroundColor Yellow

$goFiles = Get-ChildItem -Path . -Filter *.go -Recurse -File | Where-Object {
    $exclude = $false
    foreach ($dir in $excludeDirs) {
        if ($_.FullName -like "*\$dir\*") {
            $exclude = $true
            break
        }
    }
    -not $exclude
}

foreach ($file in $goFiles) {
    $totalFiles++
    try {
        $output = go fmt $file.FullName 2>&1
        if ($output) {
            $formattedFiles++
            Write-Host "  ✅ $($file.Name)" -ForegroundColor Green
        }
    }
    catch {
        Write-Host "  ❌ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "📊 gofmt 统计:" -ForegroundColor Cyan
Write-Host "  总Go文件数: $totalFiles" -ForegroundColor White
Write-Host "  已格式化文件: $formattedFiles" -ForegroundColor Green
Write-Host ""

# Step 2: go vet 检查
Write-Host "🔍 Step 2: 运行 go vet..." -ForegroundColor Yellow
Write-Host ""

# 查找所有包含go.mod的目录
$goModDirs = Get-ChildItem -Path . -Filter go.mod -Recurse -File | Where-Object {
    $exclude = $false
    foreach ($dir in $excludeDirs) {
        if ($_.FullName -like "*\$dir\*") {
            $exclude = $true
            break
        }
    }
    -not $exclude
} | ForEach-Object { $_.Directory.FullName }

$checkedModules = 0
$failedModules = 0

foreach ($dir in $goModDirs) {
    $checkedModules++
    $relativePath = $dir.Replace($PWD.Path, ".").Replace("\", "/")
    Write-Host "  检查模块: $relativePath" -ForegroundColor Cyan
    
    Push-Location $dir
    try {
        $vetOutput = go vet ./... 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "    ✅ 通过" -ForegroundColor Green
        } else {
            $failedModules++
            Write-Host "    ⚠️  发现问题:" -ForegroundColor Yellow
            $vetOutput | Where-Object { $_ -ne $null } | ForEach-Object {
                Write-Host "      $_" -ForegroundColor Yellow
            }
        }
    }
    catch {
        $failedModules++
        Write-Host "    ❌ 错误: $($_.Exception.Message)" -ForegroundColor Red
    }
    finally {
        Pop-Location
    }
}

Write-Host ""
Write-Host "📊 go vet 统计:" -ForegroundColor Cyan
Write-Host "  检查模块数: $checkedModules" -ForegroundColor White
Write-Host "  通过模块数: $($checkedModules - $failedModules)" -ForegroundColor Green
if ($failedModules -gt 0) {
    Write-Host "  失败模块数: $failedModules" -ForegroundColor Yellow
}
Write-Host ""

# Step 3: 生成报告
$reportFile = "代码格式化报告-$(Get-Date -Format 'yyyy-MM-dd-HHmm').txt"
$report = @"
代码格式化报告
生成时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')

=== gofmt 结果 ===
总Go文件数: $totalFiles
已格式化: $formattedFiles

=== go vet 结果 ===
检查模块数: $checkedModules
通过模块数: $($checkedModules - $failedModules)
失败模块数: $failedModules
"@

$report | Out-File -FilePath $reportFile -Encoding UTF8

# 最终总结
Write-Host "✅ 代码格式化完成！" -ForegroundColor Green
Write-Host "📄 详细报告已保存到: $reportFile" -ForegroundColor Cyan

if ($failedModules -gt 0) {
    Write-Host ""
    Write-Host "⚠️  警告: 有 $failedModules 个模块未通过vet检查" -ForegroundColor Yellow
    Write-Host "   请查看上面的详细输出或报告文件" -ForegroundColor Yellow
}

