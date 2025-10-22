# PowerShell Script: 生成文档统计报告
# 版本: v1.0

param(
    [string]$DocsDir = "docs",
    [string]$OutputFile = "reports/doc-statistics-$(Get-Date -Format 'yyyyMMdd').md"
)

Write-Host "📊 生成文档统计报告..." -ForegroundColor Cyan
Write-Host ""

# 统计函数
function Get-DirStats {
    param([string]$Path)
    
    $files = Get-ChildItem -Path $Path -Recurse -File -Filter "*.md"
    $totalLines = 0
    $totalWords = 0
    
    foreach ($file in $files) {
        $content = Get-Content $file.FullName -Raw -ErrorAction SilentlyContinue
        if ($content) {
            $totalLines += ($content -split "`n").Count
            $totalWords += ($content -split '\s+').Count
        }
    }
    
    return @{
        FileCount = $files.Count
        TotalLines = $totalLines
        TotalWords = $totalWords
        AvgLines = if ($files.Count -gt 0) { [Math]::Round($totalLines / $files.Count) } else { 0 }
    }
}

# 主目录统计
$mainDirs = Get-ChildItem -Path $DocsDir -Directory | Where-Object { $_.Name -notmatch "^00-" }
$allStats = @()

foreach ($dir in $mainDirs) {
    $stats = Get-DirStats -Path $dir.FullName
    $allStats += [PSCustomObject]@{
        Directory = $dir.Name
        Files = $stats.FileCount
        Lines = $stats.TotalLines
        Words = $stats.TotalWords
        AvgLines = $stats.AvgLines
    }
}

# 生成报告
$report = @"
# 📊 文档统计报告

> **生成时间**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')  
> **统计范围**: $DocsDir

---

## 总体统计

| 指标 | 数值 |
|------|-----|
| 主目录数 | $($allStats.Count) |
| 文档总数 | $(($allStats | Measure-Object -Property Files -Sum).Sum) |
| 总行数 | $(($allStats | Measure-Object -Property Lines -Sum).Sum) |
| 总字数 | $(($allStats | Measure-Object -Property Words -Sum).Sum) |

---

## 各模块统计

| 目录 | 文件数 | 总行数 | 总字数 | 平均行数 |
|------|--------|--------|--------|----------|
"@

foreach ($stat in ($allStats | Sort-Object Directory)) {
    $report += "`n| $($stat.Directory) | $($stat.Files) | $($stat.Lines) | $($stat.Words) | $($stat.AvgLines) |"
}

$report += @"


---

**生成工具**: generate_statistics.ps1  
**版本**: v1.0
"@

# 保存报告
New-Item -ItemType File -Path $OutputFile -Force | Out-Null
$report | Out-File -FilePath $OutputFile -Encoding UTF8

Write-Host "✅ 报告已生成: $OutputFile" -ForegroundColor Green
Write-Host ""
Write-Host "📊 快速预览:" -ForegroundColor Cyan
Write-Host "   总目录: $($allStats.Count)" -ForegroundColor White
Write-Host "   总文档: $(($allStats | Measure-Object -Property Files -Sum).Sum)" -ForegroundColor White
Write-Host ""

