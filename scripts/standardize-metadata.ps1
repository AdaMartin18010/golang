# 标准化所有文档的元数据格式
# 统一为：版本 | 更新日期 | 适用于

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$stats = @{
    Standardized = 0
    Scanned = 0
    Errors = 0
}

Write-Host "🔧 开始标准化元数据格式..." -ForegroundColor Cyan
Write-Host "目标格式:" -ForegroundColor Yellow
Write-Host "  **版本**: v1.0" -ForegroundColor Gray
Write-Host "  **更新日期**: 2025-10-29" -ForegroundColor Gray
Write-Host "  **适用于**: Go 1.25.3`n" -ForegroundColor Gray

$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File

foreach ($file in $files) {
    $stats.Scanned++
    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $originalContent = $content
        $modified = $false
        
        # 跳过报告文件
        if ($file.Name -match '^(🎉|🔍|📊|📝|🎯)') {
            continue
        }
        
        # 模式1: 单行格式 "**难度**: xxx | **预计阅读**: xxx | **更新**: xxx"
        if ($content -match '(?m)^(\*\*难度\*\*:.+?\|.+?\|.+?\*\*更新\*\*:.+?)$') {
            # 提取标题
            if ($content -match '^# (.+)') {
                $title = $matches[1]
                # 替换整个元数据块为标准格式
                $content = $content -replace '(?m)^(\*\*难度\*\*:.+?\|.+?\|.+?\*\*更新\*\*:.+?)$', @"
**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
"@
                $modified = $true
                if ($Verbose) {
                    Write-Host "  ✓ [单行格式] $($file.Name)" -ForegroundColor Green
                }
            }
        }
        
        # 模式2: 多行格式（各种字段混用）
        # 查找元数据块（标题后的前几行）
        if ($content -match '(?ms)^# .+?\r?\n\r?\n((?:\*\*.+?\*\*:.+?\r?\n)+)') {
            $metadataBlock = $matches[1]
            
            # 检查是否需要标准化
            $needsStandardization = $false
            if ($metadataBlock -match '\*\*(文档类型|Go版本|难度等级|最后更新|难度|预计阅读|阅读时间)\*\*:') {
                $needsStandardization = $true
            }
            
            if ($needsStandardization) {
                # 替换为标准格式
                $standardMetadata = @"
**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
"@
                
                $content = $content -replace '(?ms)(^# .+?\r?\n\r?\n)((?:\*\*.+?\*\*:.+?\r?\n)+)', "`$1$standardMetadata`r`n"
                $modified = $true
                
                if ($Verbose) {
                    Write-Host "  ✓ [多行格式] $($file.Name)" -ForegroundColor Green
                }
            }
        }
        
        # 保存修改
        if ($modified) {
            $stats.Standardized++
            
            if (-not $DryRun) {
                Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
            }
        }
        
    } catch {
        $stats.Errors++
        Write-Host "  ✗ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
    }
}

Write-Host "`n📊 统计:" -ForegroundColor Cyan
Write-Host "  扫描: $($stats.Scanned)" -ForegroundColor Gray
Write-Host "  标准化: $($stats.Standardized)" -ForegroundColor Green
Write-Host "  错误: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -gt 0){'Red'}else{'Gray'})

if ($DryRun) {
    Write-Host "`n⚠️  这是试运行模式" -ForegroundColor Yellow
}

