# 修复目录一致性问题脚本
# 功能: 确保所有Markdown文件有且只有一个目录，主题与子主题有序
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📋 修复目录一致性问题脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# GitHub Markdown anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Title)

    $anchor = $Title
    $anchor = $anchor.ToLower()
    $anchor = $anchor -replace '[^\w\s\u4e00-\u9fa5-]', ''
    $anchor = $anchor -replace '\s+', '-'
    $anchor = $anchor -replace '-+', '-'
    $anchor = $anchor.Trim('-')
    return $anchor
}

# 从Markdown文件中提取所有标题
function Get-Headings {
    param([string]$Content)

    $headings = @()
    $lines = $Content -split "`n"

    $inCodeBlock = $false

    foreach ($line in $lines) {
        $trimmed = $line.Trim()

        # 检测代码块开始和结束
        if ($trimmed -match '^```') {
            $inCodeBlock = -not $inCodeBlock
            continue
        }

        # 跳过代码块内的内容
        if ($inCodeBlock) {
            continue
        }

        # 匹配标题 (# 标题)
        if ($trimmed -match '^(#{1,6})\s+(.+)$') {
            $level = $matches[1].Length
            $title = $matches[2].Trim()

            # 跳过目录标题本身
            if ($title -match '^📋\s*目录|^目录$') {
                continue
            }

            $anchor = Get-GitHubAnchor -Title $title

            $headings += @{
                Level = $level
                Title = $title
                Anchor = $anchor
            }
        }
    }

    return $headings
}

# 生成目录
function Generate-TOC {
    param([array]$Headings)

    if ($Headings.Count -eq 0) {
        return ""
    }

    $toc = "## 📋 目录`n`n"

    # 找到第一个标题的层级作为基准
    $baseLevel = if ($Headings.Count -gt 0) { $Headings[0].Level } else { 1 }

    foreach ($heading in $Headings) {
        $level = $heading.Level
        $title = $heading.Title
        $anchor = $heading.Anchor

        # 计算缩进（相对于基准层级）
        $relativeLevel = $level - $baseLevel
        $indent = ""
        if ($relativeLevel -gt 0) {
            $indent = "  " * $relativeLevel
        }

        # 添加列表项
        $toc += "$indent- [$title](#$anchor)`n"
    }

    return $toc
}

# 检测目录位置和数量
function Get-TOCInfo {
    param([string]$Content)

    $lines = $Content -split "`n"
    $tocPositions = @()
    $inTOC = $false
    $tocStart = -1
    $tocEnd = -1

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        $trimmed = $line.Trim()

        # 检测目录标题
        if ($trimmed -match '^##+\s+.*目录') {
            if (-not $inTOC) {
                $tocStart = $i
                $inTOC = $true
                $tocPositions += @{
                    Start = $i
                    Title = $trimmed
                }
            }
        }

        # 检测目录结束（遇到下一个标题或分隔线）
        if ($inTOC) {
            if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                $tocEnd = $i - 1
                $inTOC = $false
                if ($tocPositions.Count -gt 0) {
                    $tocPositions[-1].End = $tocEnd
                }
            } elseif ($trimmed -eq '---' -and $i -gt $tocStart + 5) {
                $tocEnd = $i - 1
                $inTOC = $false
                if ($tocPositions.Count -gt 0) {
                    $tocPositions[-1].End = $tocEnd
                }
            }
        }
    }

    # 如果目录还在继续，找到结束位置
    if ($inTOC -and $tocPositions.Count -gt 0) {
        $tocPositions[-1].End = $lines.Count - 1
    }

    return $tocPositions
}

# 移除所有目录
function Remove-AllTOCs {
    param([string]$Content)

    $lines = $Content -split "`n"
    $tocInfo = Get-TOCInfo -Content $Content
    $newLines = @()
    $skipRanges = @()

    # 收集所有需要跳过的范围
    foreach ($toc in $tocInfo) {
        $skipRanges += @{
            Start = $toc.Start
            End = if ($toc.End -ge 0) { $toc.End } else { $lines.Count - 1 }
        }
    }

    # 按开始位置排序
    $skipRanges = $skipRanges | Sort-Object Start

    # 构建新内容
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $shouldSkip = $false

        foreach ($range in $skipRanges) {
            if ($i -ge $range.Start -and $i -le $range.End) {
                $shouldSkip = $true
                break
            }
        }

        if (-not $shouldSkip) {
            $newLines += $lines[$i]
        }
    }

    return $newLines -join "`n"
}

# 插入目录到正确位置
function Insert-TOC {
    param(
        [string]$Content,
        [string]$TOC
    )

    $lines = $Content -split "`n"
    $newLines = @()
    $inserted = $false
    $metadataEnd = -1

    # 查找元数据结束位置
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        $trimmed = $line.Trim()

        if ($trimmed -eq '---' -and $metadataEnd -eq -1) {
            # 检查是否是元数据后的分隔线
            $hasMetadata = $false
            for ($j = [Math]::Max(0, $i - 10); $j -lt $i; $j++) {
                if ($lines[$j] -match '\*\*版本\*\*|\*\*更新日期\*\*|\*\*适用于\*\*') {
                    $hasMetadata = $true
                    break
                }
            }

            if ($hasMetadata) {
                $metadataEnd = $i
            }
        }
    }

    # 确定插入位置（元数据后或第一个标题前）
    $insertPos = if ($metadataEnd -gt 0) { $metadataEnd + 1 } else { 0 }

    # 查找第一个实际内容标题
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()
        if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
            $insertPos = $i
            break
        }
    }

    # 构建新内容
    for ($i = 0; $i -lt $lines.Count; $i++) {
        if ($i -eq $insertPos -and -not $inserted) {
            $newLines += $TOC
            $newLines += ""
            $newLines += "---"
            $newLines += ""
            $inserted = $true
        }
        $newLines += $lines[$i]
    }

    # 如果还没插入，在末尾插入
    if (-not $inserted) {
        $newLines += ""
        $newLines += "---"
        $newLines += ""
        $newLines += $TOC
    }

    return $newLines -join "`n"
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse -File
$totalFiles = $mdFiles.Count
$processedFiles = 0
$fixedFiles = 0
$issues = @{
    NoTOC = 0
    MultipleTOC = 0
    Fixed = 0
}

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $tocInfo = Get-TOCInfo -Content $content
        $headings = Get-Headings -Content $content

        # 过滤：只保留一级和二级标题（避免目录过长）
        $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

        # 跳过标题太少的文件
        if ($filteredHeadings.Count -lt 2) {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] - $($file.Name) - 标题太少，跳过" -ForegroundColor Gray
            }
            continue
        }

        $hasIssue = $false
        $issueType = ""

        # 检查问题
        if ($tocInfo.Count -eq 0) {
            $hasIssue = $true
            $issueType = "缺少目录"
            $issues.NoTOC++
        } elseif ($tocInfo.Count -gt 1) {
            $hasIssue = $true
            $issueType = "多个目录"
            $issues.MultipleTOC++
        }

        if ($hasIssue) {
            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Yellow
            Write-Host "  问题: $issueType" -ForegroundColor Red
            Write-Host "  标题数: $($filteredHeadings.Count)" -ForegroundColor Cyan

            if (-not $DryRun) {
                # 移除所有现有目录
                $contentWithoutTOC = Remove-AllTOCs -Content $content

                # 生成新目录
                $newTOC = Generate-TOC -Headings $filteredHeadings

                # 插入新目录
                $newContent = Insert-TOC -Content $contentWithoutTOC -TOC $newTOC

                # 保存文件
                [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已修复" -ForegroundColor Green
                $issues.Fixed++
                $fixedFiles++
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 目录正常" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 修复总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  缺少目录: $($issues.NoTOC)" -ForegroundColor Yellow
Write-Host "  多个目录: $($issues.MultipleTOC)" -ForegroundColor Yellow
Write-Host "  已修复: $($issues.Fixed)" -ForegroundColor Green

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际修复" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""
