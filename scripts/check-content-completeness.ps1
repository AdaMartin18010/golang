# 内容完整性检查脚本
# 功能: 检查文档内容完整性
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$Verbose = $false
)

Write-Host "📝 内容完整性检查脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# 检查内容完整性
function Check-ContentCompleteness {
    param([string]$FilePath)

    $issues = @()
    $content = Get-Content $FilePath -Raw -Encoding UTF8

    # 检查代码块是否有说明
    $codeBlocks = [regex]::Matches($content, '```[\s\S]*?```')
    foreach ($block in $codeBlocks) {
        $beforeCode = $content.Substring([Math]::Max(0, $block.Index - 200), [Math]::Min(200, $block.Index))
        $afterCode = $content.Substring($block.Index + $block.Length, [Math]::Min(200, $content.Length - $block.Index - $block.Length))

        # 检查代码块前后是否有说明文字
        if (-not ($beforeCode -match '说明|介绍|示例|代码|实现|如下') -and
            -not ($afterCode -match '说明|介绍|示例|代码|实现|如下')) {
            $issues += @{
                Type = "代码块缺少说明"
                Line = ($content.Substring(0, $block.Index) -split "`n").Count
            }
        }
    }

    # 检查章节是否为空
    $sections = [regex]::Matches($content, '^##+\s+(.+)$', [System.Text.RegularExpressions.RegexOptions]::Multiline)
    foreach ($section in $sections) {
        $sectionTitle = $section.Groups[1].Value
        $sectionStart = $section.Index
        $nextSection = $null

        # 查找下一个同级或更高级标题
        $remaining = $content.Substring($sectionStart + $section.Length)
        $nextMatch = [regex]::Match($remaining, '^##+\s+', [System.Text.RegularExpressions.RegexOptions]::Multiline)

        if ($nextMatch.Success) {
            $sectionContent = $remaining.Substring(0, $nextMatch.Index).Trim()
        } else {
            $sectionContent = $remaining.Trim()
        }

        # 检查章节内容是否为空或过短
        if ($sectionContent.Length -lt 50 -and -not ($sectionContent -match '```|表格|列表')) {
            $issues += @{
                Type = "章节内容过短"
                Title = $sectionTitle
                Line = ($content.Substring(0, $sectionStart) -split "`n").Count
            }
        }
    }

    # 检查是否有TODO/FIXME标记
    if ($content -match 'TODO|FIXME|XXX|HACK') {
        $matches = [regex]::Matches($content, 'TODO|FIXME|XXX|HACK', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
        foreach ($match in $matches) {
            $issues += @{
                Type = "包含待办标记"
                Line = ($content.Substring(0, $match.Index) -split "`n").Count
            }
        }
    }

    return $issues
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$processedFiles = 0
$filesWithIssues = 0
$totalIssues = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

$allIssues = @()

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $issues = Check-ContentCompleteness -FilePath $file.FullName

        if ($issues.Count -gt 0) {
            $filesWithIssues++
            $totalIssues += $issues.Count

            Write-Host "[$processedFiles/$totalFiles] ⚠️ $($file.Name)" -ForegroundColor Yellow
            Write-Host "  问题数: $($issues.Count)" -ForegroundColor Red

            foreach ($issue in $issues) {
                Write-Host "  - $($issue.Type)" -ForegroundColor Gray
                if ($issue.Line) {
                    Write-Host "    行号: $($issue.Line)" -ForegroundColor DarkGray
                }
                if ($issue.Title) {
                    Write-Host "    章节: $($issue.Title)" -ForegroundColor DarkGray
                }
            }

            Write-Host ""

            $allIssues += @{
                File = $file.FullName
                Issues = $issues
            }
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 无问题" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 检查总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  有问题文件: $filesWithIssues" -ForegroundColor Yellow
Write-Host "  总问题数: $totalIssues" -ForegroundColor Red
Write-Host ""

# 按类型统计
$issueTypes = $allIssues | ForEach-Object { $_.Issues } | Group-Object -Property Type
Write-Host "📊 问题类型统计" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
foreach ($type in $issueTypes) {
    Write-Host "$($type.Name): $($type.Count) 个" -ForegroundColor Yellow
}
Write-Host ""
