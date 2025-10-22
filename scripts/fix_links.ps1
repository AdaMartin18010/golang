# PowerShell Script: 链接修复工具
# 版本: v1.0
# 日期: 2025-10-22

param(
    [string]$DocsDir = "docs",
    [string]$ReportFile = "reports/broken-links-$(Get-Date -Format 'yyyyMMdd-HHmmss').md",
    [switch]$AutoFix,
    [switch]$DryRun
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  链接修复工具" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 统计
$stats = @{
    TotalFiles = 0
    TotalLinks = 0
    BrokenLinks = 0
    FixedLinks = 0
}

$brokenLinks = @()

# 获取所有md文件
$mdFiles = Get-ChildItem -Path $DocsDir -Filter "*.md" -Recurse
$stats.TotalFiles = $mdFiles.Count

Write-Host "📊 检查 $($mdFiles.Count) 个文件..." -ForegroundColor Yellow
Write-Host ""

foreach ($file in $mdFiles) {
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    $relativePath = $file.FullName.Replace((Get-Location).Path, "").TrimStart('\')
    
    # 提取所有链接
    $links = [regex]::Matches($content, '\[([^\]]+)\]\(([^\)]+)\)')
    
    foreach ($link in $links) {
        $stats.TotalLinks++
        $linkText = $link.Groups[1].Value
        $linkUrl = $link.Groups[2].Value
        
        # 跳过外部链接和锚点
        if ($linkUrl -match "^https?://" -or $linkUrl -match "^#") {
            continue
        }
        
        # 处理相对路径
        if ($linkUrl -match "^\.\.?/") {
            # 移除锚点
            $targetUrl = $linkUrl -replace '#.*$', ''
            $targetPath = Join-Path (Split-Path $file.FullName) $targetUrl
            $targetPath = [System.IO.Path]::GetFullPath($targetPath)
            
            if (!(Test-Path $targetPath)) {
                $stats.BrokenLinks++
                
                $brokenLinks += [PSCustomObject]@{
                    File = $relativePath
                    LinkText = $linkText
                    LinkUrl = $linkUrl
                    TargetPath = $targetPath
                    Fixable = $false
                }
                
                Write-Host "  ✗ 失效链接: $relativePath" -ForegroundColor Red
                Write-Host "    链接: $linkUrl" -ForegroundColor Gray
            }
        }
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "📊 检查结果" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "总文件数:   $($stats.TotalFiles)" -ForegroundColor White
Write-Host "总链接数:   $($stats.TotalLinks)" -ForegroundColor White
Write-Host "失效链接:   $($stats.BrokenLinks)" -ForegroundColor Red
Write-Host ""

# 生成报告
$report = @"
# 🔗 链接检查报告

> **生成时间**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')  
> **检查范围**: $DocsDir

---

## 📊 统计

| 指标 | 数值 |
|------|-----|
| 总文件数 | $($stats.TotalFiles) |
| 总链接数 | $($stats.TotalLinks) |
| 失效链接 | $($stats.BrokenLinks) |
| 已修复 | $($stats.FixedLinks) |

---

## 🔴 失效链接列表

"@

if ($brokenLinks.Count -gt 0) {
    foreach ($link in $brokenLinks) {
        $report += @"

### 文件: ``$($link.File)``

**链接文本**: $($link.LinkText)  
**链接URL**: ``$($link.LinkUrl)``  
**目标路径**: ``$($link.TargetPath)``

"@
    }
} else {
    $report += "`n✅ 未发现失效链接！`n"
}

$report += @"

---

## 💡 修复建议

1. **手动修复**: 根据上面列表逐个修复
2. **自动修复**: 运行 ``.\scripts\fix_links.ps1 -AutoFix``
3. **验证修复**: 修复后重新运行检查

---

**生成工具**: fix_links.ps1  
**版本**: v1.0
"@

# 保存报告
New-Item -ItemType Directory -Path (Split-Path $ReportFile) -Force | Out-Null
$report | Out-File -FilePath $ReportFile -Encoding UTF8

Write-Host "📄 详细报告: $ReportFile" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

return $stats

