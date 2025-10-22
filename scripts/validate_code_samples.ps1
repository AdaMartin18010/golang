# scripts/validate_code_samples.ps1
# 代码示例可执行性验证工具

param (
    [string]$DocsPath = "docs-new",
    [int]$SampleLimit = 20  # 限制检查的代码示例数量
)

Write-Host "=== 📝 代码示例可执行性验证 ===" -ForegroundColor Cyan
Write-Host ""

$stats = @{
    TotalFiles = 0
    GoCodeBlocks = 0
    ValidSamples = 0
    InvalidSamples = 0
    Warnings = @()
}

# 扫描所有Markdown文件
$markdownFiles = Get-ChildItem -Path $DocsPath -Recurse -Include "*.md" | Select-Object -First 30
$stats.TotalFiles = $markdownFiles.Count

Write-Host "扫描 $($stats.TotalFiles) 个文档文件..." -ForegroundColor White
Write-Host ""

$checkedCount = 0

foreach ($file in $markdownFiles) {
    $content = Get-Content $file.FullName -Raw
    $relativePath = $file.FullName.Replace((Get-Item $DocsPath).FullName, "").TrimStart('\')
    
    # 查找Go代码块
    $codeBlocks = [regex]::Matches($content, '```go\r?\n([\s\S]*?)```')
    
    foreach ($block in $codeBlocks) {
        $stats.GoCodeBlocks++
        $code = $block.Groups[1].Value
        
        # 跳过过短的代码块（可能只是片段）
        if ($code.Length -lt 20) {
            continue
        }
        
        # 基本语法检查
        $hasPackage = $code -match '^\s*package\s+'
        $hasFunc = $code -match 'func\s+\w+\s*\('
        $hasSyntaxIssues = $false
        
        # 检查常见语法问题
        if ($code -match '\}\s*$' -and $code -notmatch '^\s*package') {
            # 代码块结束正常但没有package声明，可能是片段
            $stats.Warnings += "⚠️ 代码片段 (无package): $relativePath"
        }
        
        # 检查未闭合的括号
        $openBraces = ([regex]::Matches($code, '\{')).Count
        $closeBraces = ([regex]::Matches($code, '\}')).Count
        if ($openBraces -ne $closeBraces) {
            $hasSyntaxIssues = $true
            $stats.InvalidSamples++
            $stats.Warnings += "❌ 括号不匹配: $relativePath"
        }
        
        # 检查未闭合的引号
        $doubleQuotes = ([regex]::Matches($code, '"')).Count
        if ($doubleQuotes % 2 -ne 0) {
            $hasSyntaxIssues = $true
            $stats.InvalidSamples++
            $stats.Warnings += "❌ 引号不匹配: $relativePath"
        }
        
        if (-not $hasSyntaxIssues -and ($hasPackage -or $hasFunc)) {
            $stats.ValidSamples++
        }
        
        $checkedCount++
        if ($checkedCount -ge $SampleLimit) {
            break
        }
    }
    
    if ($checkedCount -ge $SampleLimit) {
        break
    }
}

# 生成报告
Write-Host "=== 📊 验证结果 ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "扫描统计:" -ForegroundColor Yellow
Write-Host "  • 检查文件: $($stats.TotalFiles)" -ForegroundColor White
Write-Host "  • Go代码块: $($stats.GoCodeBlocks)" -ForegroundColor White
Write-Host "  • 有效示例: $($stats.ValidSamples)" -ForegroundColor Green
Write-Host "  • 语法问题: $($stats.InvalidSamples)" -ForegroundColor $(if ($stats.InvalidSamples -eq 0) { "Green" } else { "Red" })

if ($stats.Warnings.Count -gt 0 -and $stats.Warnings.Count -le 10) {
    Write-Host ""
    Write-Host "⚠️ 问题详情:" -ForegroundColor Yellow
    $stats.Warnings | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
} elseif ($stats.Warnings.Count -gt 10) {
    Write-Host ""
    Write-Host "⚠️ 问题详情 (显示前10条，共$($stats.Warnings.Count)条):" -ForegroundColor Yellow
    $stats.Warnings | Select-Object -First 10 | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
}

# 计算评分
if ($stats.GoCodeBlocks -gt 0) {
    $successRate = [Math]::Round(($stats.ValidSamples * 100.0 / $stats.GoCodeBlocks), 1)
} else {
    $successRate = 100
}

Write-Host ""
Write-Host "🎯 代码质量: $successRate%" -ForegroundColor $(
    if ($successRate -ge 90) { "Green" }
    elseif ($successRate -ge 75) { "Yellow" }
    else { "Red" }
)

Write-Host ""
Write-Host "✅ 代码验证完成！" -ForegroundColor Green
Write-Host "   注：本工具进行基础语法检查，实际可执行性需要完整的编译测试。" -ForegroundColor DarkGray

# 保存报告
$reportPath = "reports/Phase6-代码验证报告-$(Get-Date -Format 'yyyy-MM-dd').md"
$reportContent = @"
# Phase 6 - 代码示例验证报告

**生成时间**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')  
**检查范围**: $DocsPath  
**代码质量**: $successRate%

---

## 📊 验证统计

| 项目 | 数量 | 状态 |
|------|------|------|
| 检查文件 | $($stats.TotalFiles) | ✅ |
| Go代码块 | $($stats.GoCodeBlocks) | ✅ |
| 有效示例 | $($stats.ValidSamples) | ✅ |
| 语法问题 | $($stats.InvalidSamples) | $(if ($stats.InvalidSamples -eq 0) { "✅" } else { "⚠️" }) |

---

## 🎯 质量评估

- **成功率**: $successRate%
- **质量等级**: $(
    if ($successRate -ge 90) { "A (优秀)" }
    elseif ($successRate -ge 75) { "B (良好)" }
    else { "C (需改进)" }
)

---

## ⚠️ 发现的问题

$(if ($stats.Warnings.Count -eq 0) {
    "✅ 未发现语法问题"
} else {
    $stats.Warnings -join "`n"
})

---

## 💡 建议

1. **代码片段**: 文档中包含大量教学性代码片段，这是正常的
2. **完整性**: 建议为关键示例提供完整的可运行代码
3. **测试**: 对于教程文档，建议建立CI自动测试代码示例

---

**维护者**: Documentation Team  
**创建日期**: $(Get-Date -Format 'yyyy-MM-dd')  
**文档状态**: ✅ 完成
"@

Set-Content -Path $reportPath -Value $reportContent -Encoding UTF8
Write-Host "📄 详细报告已保存: $reportPath" -ForegroundColor Cyan

