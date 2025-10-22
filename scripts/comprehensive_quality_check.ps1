# scripts/comprehensive_quality_check.ps1
# 全面质量检查工具 - Phase 6

param (
    [string]$DocsPath = "docs-new"
)

Write-Host "=== 📋 全面质量检查工具 ===" -ForegroundColor Cyan
Write-Host ""

# 初始化统计
$stats = @{
    TotalFiles = 0
    BrokenLinks = 0
    FormattingIssues = 0
    MissingMetadata = 0
    CodeBlockIssues = 0
    Warnings = @()
    Errors = @()
}

# 1. 检查所有Markdown文件
Write-Host "1️⃣ 扫描文档文件..." -ForegroundColor Yellow
$markdownFiles = Get-ChildItem -Path $DocsPath -Recurse -Include "*.md"
$stats.TotalFiles = $markdownFiles.Count
Write-Host "   发现 $($stats.TotalFiles) 个文档文件" -ForegroundColor White

# 2. 链接验证
Write-Host ""
Write-Host "2️⃣ 验证内部链接..." -ForegroundColor Yellow
foreach ($file in $markdownFiles) {
    $content = Get-Content $file.FullName -Raw
    $relativePath = $file.FullName.Replace((Get-Item $DocsPath).FullName, "").TrimStart('\')
    
    # 查找所有Markdown链接 [文本](路径)
    $links = [regex]::Matches($content, '\[([^\]]+)\]\(([^\)]+)\)')
    
    foreach ($link in $links) {
        $linkText = $link.Groups[1].Value
        $linkPath = $link.Groups[2].Value
        
        # 跳过外部链接和锚点链接
        if ($linkPath -match '^https?://' -or $linkPath -match '^#') {
            continue
        }
        
        # 分离文件路径和锚点
        $filePart = $linkPath
        $anchorPart = ""
        if ($linkPath -match '(.+)(#.+)$') {
            $filePart = $matches[1]
            $anchorPart = $matches[2]
        }
        
        # 构建完整路径
        $targetPath = $filePart
        if (-not [System.IO.Path]::IsPathRooted($filePart)) {
            $targetPath = Join-Path (Split-Path $file.FullName) $filePart
            $targetPath = [System.IO.Path]::GetFullPath($targetPath)
        }
        
        # 检查文件是否存在（只验证文件路径部分，不验证锚点）
        if (-not (Test-Path $targetPath)) {
            $stats.BrokenLinks++
            $stats.Errors += "🔗 失效链接: $relativePath -> $linkPath"
        }
    }
}
Write-Host "   发现 $($stats.BrokenLinks) 个失效链接" -ForegroundColor $(if ($stats.BrokenLinks -eq 0) { "Green" } else { "Red" })

# 3. 格式检查
Write-Host ""
Write-Host "3️⃣ 检查文档格式..." -ForegroundColor Yellow
foreach ($file in $markdownFiles) {
    $content = Get-Content $file.FullName -Raw
    $lines = Get-Content $file.FullName
    $relativePath = $file.FullName.Replace((Get-Item $DocsPath).FullName, "").TrimStart('\')
    
    # 检查H1标题
    if ($file.Name -eq "README.md") {
        if ($content -notmatch '^\s*#\s+\S') {
            $stats.FormattingIssues++
            $stats.Warnings += "⚠️ 缺少H1标题: $relativePath"
        }
    }
    
    # 检查目录结构
    if ($content -match '##\s+目录' -or $content -match '##\s+Table of Contents') {
        # 有目录，检查是否在正确位置（应该在元信息之后）
        $tocIndex = $content.IndexOf('## 目录')
        if ($tocIndex -eq -1) {
            $tocIndex = $content.IndexOf('## Table of Contents')
        }
        
        # 检查目录项是否完整
        $h2Headers = [regex]::Matches($content, '##\s+([^#\n]+)')
        $tocLinks = [regex]::Matches($content, '\[([^\]]+)\]\(#[^\)]+\)')
        
        if ($h2Headers.Count -gt $tocLinks.Count + 2) {  # +2 for metadata and TOC itself
            $stats.Warnings += "⚠️ 目录不完整: $relativePath (有 $($h2Headers.Count) 个章节，但只有 $($tocLinks.Count) 个目录项)"
        }
    }
    
    # 检查代码块是否闭合
    $codeBlockCount = ([regex]::Matches($content, '```')).Count
    if ($codeBlockCount % 2 -ne 0) {
        $stats.CodeBlockIssues++
        $stats.Errors += "❌ 代码块未闭合: $relativePath"
    }
}
Write-Host "   发现 $($stats.FormattingIssues) 个格式问题" -ForegroundColor $(if ($stats.FormattingIssues -eq 0) { "Green" } else { "Yellow" })
Write-Host "   发现 $($stats.CodeBlockIssues) 个代码块问题" -ForegroundColor $(if ($stats.CodeBlockIssues -eq 0) { "Green" } else { "Red" })

# 4. 元信息检查
Write-Host ""
Write-Host "4️⃣ 检查元信息..." -ForegroundColor Yellow
foreach ($file in $markdownFiles) {
    $content = Get-Content $file.FullName -Raw
    $relativePath = $file.FullName.Replace((Get-Item $DocsPath).FullName, "").TrimStart('\')
    
    # 检查是否有维护者信息
    if ($content -notmatch '维护者|创建日期|最后更新|文档状态') {
        $stats.MissingMetadata++
        $stats.Warnings += "⚠️ 缺少元信息: $relativePath"
    }
}
Write-Host "   发现 $($stats.MissingMetadata) 个文件缺少元信息" -ForegroundColor $(if ($stats.MissingMetadata -eq 0) { "Green" } else { "Yellow" })

# 5. 生成报告
Write-Host ""
Write-Host "=== 📊 质量检查报告 ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "📈 统计汇总:" -ForegroundColor Yellow
Write-Host "  • 检查文件: $($stats.TotalFiles)" -ForegroundColor White
Write-Host "  • 失效链接: $($stats.BrokenLinks)" -ForegroundColor $(if ($stats.BrokenLinks -eq 0) { "Green" } else { "Red" })
Write-Host "  • 格式问题: $($stats.FormattingIssues)" -ForegroundColor $(if ($stats.FormattingIssues -eq 0) { "Green" } else { "Yellow" })
Write-Host "  • 代码块问题: $($stats.CodeBlockIssues)" -ForegroundColor $(if ($stats.CodeBlockIssues -eq 0) { "Green" } else { "Red" })
Write-Host "  • 缺少元信息: $($stats.MissingMetadata)" -ForegroundColor $(if ($stats.MissingMetadata -eq 0) { "Green" } else { "Yellow" })

if ($stats.Errors.Count -gt 0) {
    Write-Host ""
    Write-Host "❌ 错误详情 ($($stats.Errors.Count)):" -ForegroundColor Red
    $stats.Errors | ForEach-Object { Write-Host "  $_" -ForegroundColor Red }
}

if ($stats.Warnings.Count -gt 0 -and $stats.Warnings.Count -le 20) {
    Write-Host ""
    Write-Host "⚠️ 警告详情 ($($stats.Warnings.Count)):" -ForegroundColor Yellow
    $stats.Warnings | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
} elseif ($stats.Warnings.Count -gt 20) {
    Write-Host ""
    Write-Host "⚠️ 警告详情 (显示前20条，共$($stats.Warnings.Count)条):" -ForegroundColor Yellow
    $stats.Warnings | Select-Object -First 20 | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
}

# 6. 质量评分
$totalIssues = $stats.BrokenLinks + $stats.FormattingIssues + $stats.CodeBlockIssues + $stats.MissingMetadata
$qualityScore = [Math]::Max(0, 100 - ($totalIssues * 100.0 / $stats.TotalFiles))

Write-Host ""
Write-Host "🏆 质量评分: $([Math]::Round($qualityScore, 1))/100" -ForegroundColor $(
    if ($qualityScore -ge 95) { "Green" }
    elseif ($qualityScore -ge 85) { "Yellow" }
    else { "Red" }
)

if ($qualityScore -ge 95) {
    Write-Host "   等级: A+ (卓越)" -ForegroundColor Green
} elseif ($qualityScore -ge 90) {
    Write-Host "   等级: A (优秀)" -ForegroundColor Green
} elseif ($qualityScore -ge 85) {
    Write-Host "   等级: B+ (良好)" -ForegroundColor Yellow
} elseif ($qualityScore -ge 80) {
    Write-Host "   等级: B (合格)" -ForegroundColor Yellow
} else {
    Write-Host "   等级: C (需改进)" -ForegroundColor Red
}

Write-Host ""
Write-Host "✅ 质量检查完成！" -ForegroundColor Green

# 7. 保存详细报告
$reportPath = "reports/Phase6-质量检查报告-$(Get-Date -Format 'yyyy-MM-dd').md"
$reportContent = @"
# Phase 6 - 全面质量检查报告

**生成时间**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')  
**检查范围**: $DocsPath  
**质量评分**: $([Math]::Round($qualityScore, 1))/100

---

## 📊 统计汇总

| 项目 | 数量 | 状态 |
|------|------|------|
| 检查文件总数 | $($stats.TotalFiles) | ✅ |
| 失效链接 | $($stats.BrokenLinks) | $(if ($stats.BrokenLinks -eq 0) { "✅" } else { "❌" }) |
| 格式问题 | $($stats.FormattingIssues) | $(if ($stats.FormattingIssues -eq 0) { "✅" } else { "⚠️" }) |
| 代码块问题 | $($stats.CodeBlockIssues) | $(if ($stats.CodeBlockIssues -eq 0) { "✅" } else { "❌" }) |
| 缺少元信息 | $($stats.MissingMetadata) | $(if ($stats.MissingMetadata -eq 0) { "✅" } else { "⚠️" }) |

---

## 🏆 质量评级

- **总体评分**: $([Math]::Round($qualityScore, 1))/100
- **质量等级**: $(
    if ($qualityScore -ge 95) { "A+ (卓越)" }
    elseif ($qualityScore -ge 90) { "A (优秀)" }
    elseif ($qualityScore -ge 85) { "B+ (良好)" }
    elseif ($qualityScore -ge 80) { "B (合格)" }
    else { "C (需改进)" }
)

---

## ❌ 错误详情

$(if ($stats.Errors.Count -eq 0) {
    "✅ 无错误发现"
} else {
    $stats.Errors -join "`n"
})

---

## ⚠️ 警告详情

$(if ($stats.Warnings.Count -eq 0) {
    "✅ 无警告"
} else {
    $stats.Warnings -join "`n"
})

---

## 📋 建议措施

$(if ($totalIssues -eq 0) {
    "✅ 文档质量优秀，无需额外改进措施。"
} else {
    $suggestions = @()
    if ($stats.BrokenLinks -gt 0) {
        $suggestions += "1. 修复失效链接，确保所有内部链接指向有效文档"
    }
    if ($stats.FormattingIssues -gt 0) {
        $suggestions += "2. 统一文档格式，确保所有README都有正确的H1标题"
    }
    if ($stats.CodeBlockIssues -gt 0) {
        $suggestions += "3. 修复未闭合的代码块"
    }
    if ($stats.MissingMetadata -gt 0) {
        $suggestions += "4. 补充缺失的元信息（维护者、日期、状态）"
    }
    $suggestions -join "`n"
})

---

**维护者**: Documentation Team  
**创建日期**: $(Get-Date -Format 'yyyy-MM-dd')  
**文档状态**: ✅ 完成
"@

Set-Content -Path $reportPath -Value $reportContent -Encoding UTF8
Write-Host "📄 详细报告已保存: $reportPath" -ForegroundColor Cyan

