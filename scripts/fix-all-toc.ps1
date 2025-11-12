# 全面修复目录脚本
# 功能: 确保所有Markdown文件有且只有一个目录，目录有序且完整
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📋 全面修复目录脚本" -ForegroundColor Cyan
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

    foreach ($line in $lines) {
        $trimmed = $line.Trim()
        if ($trimmed -match '^(#{1,6})\s+(.+)$') {
            $level = $matches[1].Length
            $title = $matches[2].Trim()
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
    $prevLevel = 0

    foreach ($heading in $Headings) {
        $level = $heading.Level
        $title = $heading.Title
        $anchor = $heading.Anchor

        # 只包含一级和二级标题（避免目录过长）
        if ($level -gt 2) {
            continue
        }

        # 计算缩进
        $indent = ""
        if ($level -gt 1) {
            $indent = "  " * ($level - 1)
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
    $tocStart = -1
    $tocEnd = -1

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()

        # 检测目录标题
        if ($trimmed -match '^##+\s+.*目录') {
            if ($tocStart -eq -1) {
                $tocStart = $i
            }
            $tocPositions += $i
        }

        # 检测目录结束（遇到下一个标题）
        if ($tocStart -ge 0 -and $tocEnd -eq -1) {
            if ($i -gt $tocStart -and $trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                $tocEnd = $i
                break
            }
        }
    }

    return @{
        Count = $tocPositions.Count
        Positions = $tocPositions
        Start = $tocStart
        End = $tocEnd
    }
}

# 修复目录
function Fix-TOC {
    param([string]$Content)

    $lines = $Content -split "`n"
    $tocInfo = Get-TOCInfo -Content $Content

    # 提取标题
    $headings = Get-Headings -Content $Content

    # 过滤：只保留一级和二级标题
    $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

    # 如果标题少于2个，不需要目录
    if ($filteredHeadings.Count -lt 2) {
        return $Content
    }

    # 生成新目录
    $newTOC = Generate-TOC -Headings $filteredHeadings

    # 找到元数据结束位置
    $metadataEnd = -1
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()
        if ($trimmed -eq '---' -and $i -gt 0) {
            $metadataEnd = $i
            break
        }
    }

    # 确定插入位置（元数据后或第一个标题前）
    $insertPos = if ($metadataEnd -gt 0) { $metadataEnd + 1 } else { 0 }

    # 找到第一个标题位置
    $firstHeading = -1
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()
        if ($trimmed -match '^#\s+') {
            $firstHeading = $i
            break
        }
    }

    if ($firstHeading -gt 0 -and $insertPos -gt $firstHeading) {
        $insertPos = $firstHeading + 1
    }

    # 如果有多个目录，删除所有旧目录
    $newLines = @()
    $tocInserted = $false

    for ($i = 0; $i -lt $lines.Count; $i++) {
        # 跳过旧的目录区域
        if ($tocInfo.Start -ge 0 -and $i -ge $tocInfo.Start) {
            if ($tocInfo.End -gt 0 -and $i -lt $tocInfo.End) {
                continue
            } elseif ($tocInfo.End -eq -1) {
                # 如果没找到结束位置，跳过直到下一个标题
                $trimmed = $lines[$i].Trim()
                if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                    # 遇到下一个标题，停止跳过
                } else {
                    continue
                }
            }
        }

        # 在插入位置添加新目录（只添加一次）
        if (-not $tocInserted -and $i -eq $insertPos) {
            $newLines += $newTOC.TrimEnd()
            $newLines += ""
            $newLines += "---"
            $newLines += ""
            $tocInserted = $true
        }

        $newLines += $lines[$i]
    }

    # 如果插入位置在文件末尾，在末尾添加
    if (-not $tocInserted -and $insertPos -ge $lines.Count) {
        $newLines += ""
        $newLines += $newTOC.TrimEnd()
        $newLines += ""
        $newLines += "---"
        $newLines += ""
    }

    return $newLines -join "`n"
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse -File
$totalFiles = $mdFiles.Count
$processedFiles = 0
$fixedFiles = 0
$issues = @()

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $tocInfo = Get-TOCInfo -Content $content

        # 提取标题
        $headings = Get-Headings -Content $content
        $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

        # 如果标题少于2个，跳过
        if ($filteredHeadings.Count -lt 2) {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] - $($file.Name) - 标题太少，跳过" -ForegroundColor Gray
            }
            continue
        }

        $hasIssue = $false
        $issueMsg = ""

        # 检查问题
        if ($tocInfo.Count -eq 0) {
            $hasIssue = $true
            $issueMsg = "缺少目录"
        } elseif ($tocInfo.Count -gt 1) {
            $hasIssue = $true
            $issueMsg = "有 $($tocInfo.Count) 个目录（应只有1个）"
        }

        if ($hasIssue) {
            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Cyan
            Write-Host "  问题: $issueMsg" -ForegroundColor Yellow
            Write-Host "  标题数: $($filteredHeadings.Count)" -ForegroundColor Gray

            if (-not $DryRun) {
                $newContent = Fix-TOC -Content $content
                if ($newContent -ne $content) {
                    [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                    Write-Host "  ✅ 已修复" -ForegroundColor Green
                    $fixedFiles++
                }
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }

            $issues += @{
                File = $file.FullName.Replace((Get-Location).Path + "\", "")
                Issue = $issueMsg
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
Write-Host "  发现问题: $($issues.Count)" -ForegroundColor Yellow
Write-Host "  已修复: $fixedFiles" -ForegroundColor Green

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际修复" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

if ($issues.Count -gt 0) {
    Write-Host ""
    Write-Host "问题文件列表:" -ForegroundColor Yellow
    $issues | ForEach-Object {
        Write-Host "  - $($_.File): $($_.Issue)" -ForegroundColor Gray
    }
}

Write-Host ""
