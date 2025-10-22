# PowerShell Script: 文档质量检查
# 版本: v1.0
# 日期: 2025-10-22

param(
    [string]$DocsDir = "docs",
    [string]$OutputFile = "reports/quality-check-$(Get-Date -Format 'yyyyMMdd-HHmmss').md"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  文档质量检查脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$issues = @()
$totalFiles = 0
$checkedFiles = 0

# 检查所有.md文件
$mdFiles = Get-ChildItem -Path $DocsDir -Filter "*.md" -Recurse

Write-Host "📊 开始检查 $($mdFiles.Count) 个文件..." -ForegroundColor Yellow
Write-Host ""

foreach ($file in $mdFiles) {
    $totalFiles++
    $relativePath = $file.FullName.Replace((Get-Location).Path, "").TrimStart('\')
    
    Write-Host "  检查: $relativePath" -ForegroundColor Gray
    
    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $checkedFiles++
        
        # 检查1: 元信息
        if ($content -notmatch "维护者:|最后更新:|文档状态:") {
            $issues += [PSCustomObject]@{
                File = $relativePath
                Type = "元信息缺失"
                Severity = "中"
                Description = "缺少必要的元信息（维护者/更新时间/状态）"
            }
        }
        
        # 检查2: 标题
        if ($content -notmatch "^# ") {
            $issues += [PSCustomObject]@{
                File = $relativePath
                Type = "缺少标题"
                Severity = "高"
                Description = "文档缺少一级标题"
            }
        }
        
        # 检查3: 链接
        $links = [regex]::Matches($content, '\[([^\]]+)\]\(([^\)]+)\)')
        foreach ($link in $links) {
            $url = $link.Groups[2].Value
            if ($url -match "^\.\.?/") {
                $targetPath = Join-Path (Split-Path $file.FullName) $url
                $targetPath = $targetPath -replace '#.*$', ''  # 移除锚点
                if (!(Test-Path $targetPath)) {
                    $issues += [PSCustomObject]@{
                        File = $relativePath
                        Type = "失效链接"
                        Severity = "高"
                        Description = "链接目标不存在: $url"
                    }
                }
            }
        }
        
        # 检查4: 代码块
        $codeBlocks = [regex]::Matches($content, '```(\w+)\r?\n(.*?)```', [System.Text.RegularExpressions.RegexOptions]::Singleline)
        foreach ($block in $codeBlocks) {
            $lang = $block.Groups[1].Value
            $code = $block.Groups[2].Value
            
            if ($lang -eq "go" -and $code -notmatch "package\s+\w+") {
                $issues += [PSCustomObject]@{
                    File = $relativePath
                    Type = "代码不完整"
                    Severity = "低"
                    Description = "Go代码缺少package声明"
                }
            }
        }
        
        # 检查5: 文件大小
        $lines = ($content -split "`n").Count
        if ($lines -lt 50) {
            $issues += [PSCustomObject]@{
                File = $relativePath
                Type = "内容过少"
                Severity = "低"
                Description = "文档行数过少（<50行），可能需要扩充"
            }
        }
        
    } catch {
        Write-Host "  ⚠️  读取失败: $_" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan

# 生成报告
$report = @"
# 📊 文档质量检查报告

> **生成时间**: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')  
> **检查范围**: $DocsDir  
> **检查文件**: $checkedFiles / $totalFiles

---

## 总体统计

| 指标 | 数值 |
|------|-----|
| 总文件数 | $totalFiles |
| 已检查 | $checkedFiles |
| 发现问题 | $($issues.Count) |
| 高严重性 | $(($issues | Where-Object Severity -eq '高').Count) |
| 中严重性 | $(($issues | Where-Object Severity -eq '中').Count) |
| 低严重性 | $(($issues | Where-Object Severity -eq '低').Count) |

---

## 问题列表

### 🔴 高严重性问题

"@

$highIssues = $issues | Where-Object Severity -eq '高'
if ($highIssues.Count -gt 0) {
    foreach ($issue in $highIssues) {
        $report += @"

**文件**: ``$($issue.File)``  
**类型**: $($issue.Type)  
**说明**: $($issue.Description)

"@
    }
} else {
    $report += "`n无高严重性问题`n"
}

$report += @"

### 🟡 中严重性问题

"@

$mediumIssues = $issues | Where-Object Severity -eq '中'
if ($mediumIssues.Count -gt 0) {
    foreach ($issue in $mediumIssues | Select-Object -First 10) {
        $report += @"

**文件**: ``$($issue.File)``  
**类型**: $($issue.Type)  
**说明**: $($issue.Description)

"@
    }
    if ($mediumIssues.Count -gt 10) {
        $report += "`n*...还有 $(($mediumIssues.Count - 10)) 个中严重性问题*`n"
    }
} else {
    $report += "`n无中严重性问题`n"
}

$report += @"

### 🟢 低严重性问题

总数: $(($issues | Where-Object Severity -eq '低').Count) 个（详见完整日志）

---

## 建议

"@

if ($issues.Count -eq 0) {
    $report += "✅ 未发现问题，文档质量良好！`n"
} else {
    $report += @"
1. 优先修复高严重性问题
2. 补充缺失的元信息
3. 修复失效链接
4. 完善代码示例
5. 扩充内容过少的文档

---

## 下一步

- [ ] 修复所有高严重性问题
- [ ] 修复中严重性问题
- [ ] 运行 lychee 进行深度链接检查
- [ ] 运行 prettier 进行格式检查

"@
}

$report += @"

---

**生成工具**: check_quality.ps1  
**版本**: v1.0
"@

# 保存报告
New-Item -ItemType File -Path $OutputFile -Force | Out-Null
$report | Out-File -FilePath $OutputFile -Encoding UTF8

# 控制台输出
if ($issues.Count -eq 0) {
    Write-Host "✅ 质量检查通过! 未发现问题" -ForegroundColor Green
} else {
    Write-Host "⚠️  发现 $($issues.Count) 个问题:" -ForegroundColor Yellow
    Write-Host "   🔴 高严重性: $(($issues | Where-Object Severity -eq '高').Count)" -ForegroundColor Red
    Write-Host "   🟡 中严重性: $(($issues | Where-Object Severity -eq '中').Count)" -ForegroundColor Yellow
    Write-Host "   🟢 低严重性: $(($issues | Where-Object Severity -eq '低').Count)" -ForegroundColor Green
}

Write-Host ""
Write-Host "📄 详细报告: $OutputFile" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

