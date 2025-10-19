# 🔧 版本号批量替换脚本
# 将所有 Go 1.25 引用替换为 Go 1.23+

Write-Host "🔧 开始批量替换版本号..." -ForegroundColor Cyan
Write-Host ""

# 统计
$totalFiles = 0
$replacedFiles = 0

# 定义要搜索的文件类型
$fileTypes = @("*.md", "*.go", "*.txt")

# 定义替换规则
$replacements = @{
    'Go 1\.25\.3' = 'Go 1.23+'
    'Go 1\.25\.2' = 'Go 1.23+'
    'Go 1\.25\.1' = 'Go 1.23+'
    'Go 1\.25' = 'Go 1.23+'
    'go1\.25' = 'go1.23'
    'go version go1\.25' = 'go version go1.23'
    '1\.25\.3' = '1.23+'
    '1\.25\.2' = '1.23+'
    '1\.25\.1' = '1.23+'
    'Go-1\.25' = 'Go-1.23'
    'Go 1.25新特性' = 'Go 1.23+现代特性'
    'Go 1.25特性' = 'Go 1.23+特性'
    'Go 1.25的' = 'Go 1.23+的'
}

# 排除的目录
$excludeDirs = @(
    ".git",
    "node_modules",
    ".cursor",
    "vendor"
)

foreach ($fileType in $fileTypes) {
    Write-Host "📝 处理 $fileType 文件..." -ForegroundColor Yellow
    
    $files = Get-ChildItem -Path . -Filter $fileType -Recurse -File | Where-Object {
        $exclude = $false
        foreach ($dir in $excludeDirs) {
            if ($_.FullName -like "*\$dir\*") {
                $exclude = $true
                break
            }
        }
        -not $exclude
    }
    
    foreach ($file in $files) {
        $totalFiles++
        $content = Get-Content $file.FullName -Raw -Encoding UTF8 -ErrorAction SilentlyContinue
        
        if ($null -eq $content) {
            continue
        }
        
        $originalContent = $content
        $changed = $false
        
        # 应用所有替换规则
        foreach ($pattern in $replacements.Keys) {
            $replacement = $replacements[$pattern]
            if ($content -match $pattern) {
                $content = $content -replace $pattern, $replacement
                $changed = $true
            }
        }
        
        if ($changed) {
            try {
                Set-Content $file.FullName $content -Encoding UTF8 -NoNewline
                $replacedFiles++
                Write-Host "  ✅ $($file.Name)" -ForegroundColor Green
            }
            catch {
                Write-Host "  ❌ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
            }
        }
    }
}

Write-Host ""
Write-Host "📊 替换统计:" -ForegroundColor Cyan
Write-Host "  总文件数: $totalFiles" -ForegroundColor White
Write-Host "  已替换文件: $replacedFiles" -ForegroundColor Green
Write-Host ""

if ($replacedFiles -gt 0) {
    Write-Host "✅ 版本号替换完成！" -ForegroundColor Green
} else {
    Write-Host "ℹ️  没有文件需要替换" -ForegroundColor Yellow
}

# 生成替换报告
$reportFile = "版本替换报告-$(Get-Date -Format 'yyyy-MM-dd-HHmm').txt"
$report = @"
版本号替换报告
生成时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')

处理文件数: $totalFiles
替换文件数: $replacedFiles

替换规则:
$(($replacements.GetEnumerator() | ForEach-Object { "  $($_.Key) → $($_.Value)" }) -join "`n")
"@

$report | Out-File -FilePath $reportFile -Encoding UTF8
Write-Host "📄 详细报告已保存到: $reportFile" -ForegroundColor Cyan

