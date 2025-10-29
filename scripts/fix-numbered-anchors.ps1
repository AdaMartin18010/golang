# 修复数字序号anchor问题
# 例如：[1.1 xxx](#1.1-xxx) -> [1.1 xxx](#11-xxx)

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    FilesProcessed = 0
    LinksFixed = 0
    Errors = 0
}

Write-Host "🔢 修复数字序号anchor...`n" -ForegroundColor Cyan

# 处理文件
$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File

foreach ($file in $files) {
    try {
        $stats.FilesProcessed++
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        
        if ([string]::IsNullOrWhiteSpace($content)) {
            continue
        }
        
        $originalContent = $content
        $fixCount = 0
        
        # 模式1: (#1.1-xxx) -> (#11-xxx) 
        # 匹配 (#数字.数字-) 改为 (#数字数字-)
        $pattern1 = '\(#(\d+)\.(\d+)-'
        $replacement1 = '(#$1$2-'
        if ($content -match $pattern1) {
            $content = $content -replace $pattern1, $replacement1
            $count = ([regex]::Matches($originalContent, $pattern1)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复模式1 (#X.Y-): $count 个" -ForegroundColor Gray
            }
        }
        
        # 模式2: (#10.1-xxx) -> (#101-xxx) 
        # 三位数也要处理
        $pattern2 = '\(#(\d+)\.(\d+)\.(\d+)-'
        $replacement2 = '(#$1$2$3-'
        if ($content -match $pattern2) {
            $content = $content -replace $pattern2, $replacement2
            $count = ([regex]::Matches($originalContent, $pattern2)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复模式2 (#X.Y.Z-): $count 个" -ForegroundColor Gray
            }
        }
        
        # 模式3: (#1.1) -> (#11) (末尾的情况)
        $pattern3 = '\(#(\d+)\.(\d+)\)'
        $replacement3 = '(#$1$2)'
        # 确保不是已经修复过的
        if ($content -match $pattern3 -and $content -notmatch '\(#\d{3,}\)') {
            $content = $content -replace $pattern3, $replacement3
            $count = ([regex]::Matches($originalContent, $pattern3)).Count
            $fixCount += $count
            if ($Verbose) {
                Write-Host "  修复模式3 (#X.Y): $count 个" -ForegroundColor Gray
            }
        }
        
        if ($fixCount -gt 0) {
            Write-Host "✓ $($file.FullName -replace [regex]::Escape($PWD), '.')" -ForegroundColor Green
            Write-Host "  修复: $fixCount 个链接" -ForegroundColor Gray
            
            $stats.LinksFixed += $fixCount
            
            if (-not $DryRun) {
                Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
            }
        }
        
    } catch {
        $stats.Errors++
        Write-Host "✗ $($file.Name): $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 显示统计
Write-Host "`n" + ("="*70) -ForegroundColor Cyan
Write-Host "📊 修复统计:" -ForegroundColor Cyan
Write-Host "  处理文件: $($stats.FilesProcessed)" -ForegroundColor White
Write-Host "  修复链接: $($stats.LinksFixed)" -ForegroundColor Green
Write-Host "  错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -eq 0){'Green'}else{'Red'})
Write-Host ("="*70) -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "`n⚠️  试运行模式 - 未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复完成！" -ForegroundColor Green
}

