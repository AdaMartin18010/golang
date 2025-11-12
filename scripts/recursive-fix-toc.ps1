# 全面递归修复目录脚本
# 功能: 递归检查并修复所有Markdown文件的目录问题
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📋 全面递归修复目录脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# GitHub Markdown anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Title)
    $anchor = $Title.ToLower()
    $anchor = $anchor -replace '[^\w\s\u4e00-\u9fa5-]', ''
    $anchor = $anchor -replace '\s+', '-'
    $anchor = $anchor -replace '-+', '-'
    $anchor = $anchor.Trim('-')
    return $anchor
}

# 提取标题
function Get-Headings {
    param([string]$Content)
    $headings = @()
    $lines = $Content -split "`n"
    foreach ($line in $lines) {
        $trimmed = $line.Trim()
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
    if ($Headings.Count -eq 0) { return "" }

    $toc = "## 📋 目录`n`n"
    foreach ($heading in $Headings) {
        $level = $heading.Level
        # 只包含一级和二级标题
        if ($level -gt 2) { continue }
        $indent = if ($level -gt 1) { "  " * ($level - 1) } else { "" }
        $toc += "$indent- [$($heading.Title)](#$($heading.Anchor))`n"
    }
    return $toc
}

# 查找所有目录位置
function Find-TOCRanges {
    param([array]$Lines)
    $ranges = @()
    $tocStart = -1

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()

        # 检测目录标题
        if ($trimmed -match '^##+\s+.*目录') {
            if ($tocStart -eq -1) {
                $tocStart = $i
            }
        } elseif ($tocStart -ge 0) {
            # 查找目录结束位置
            if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                # 遇到下一个标题，目录结束
                $ranges += @{ Start = $tocStart; End = $i }
                $tocStart = -1
            } elseif ($trimmed -eq '---' -and $i -gt $tocStart + 5) {
                # 遇到分隔线且距离目录标题较远，可能是目录结束
                $ranges += @{ Start = $tocStart; End = $i }
                $tocStart = -1
            } elseif ($i -eq $lines.Count - 1) {
                # 文件末尾
                $ranges += @{ Start = $tocStart; End = $i + 1 }
                $tocStart = -1
            }
        }
    }

    # 如果还有未结束的目录
    if ($tocStart -ge 0) {
        $ranges += @{ Start = $tocStart; End = $lines.Count }
    }

    return $ranges
}

# 修复目录
function Fix-TOC {
    param([string]$Content)

    $lines = $Content -split "`n"
    $headings = Get-Headings -Content $Content
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

    # 找到第一个标题位置
    $firstHeading = -1
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()
        if ($trimmed -match '^#\s+') {
            $firstHeading = $i
            break
        }
    }

    # 确定插入位置（元数据后，第一个标题前）
    $insertPos = if ($metadataEnd -gt 0) {
        $metadataEnd + 1
    } else {
        if ($firstHeading -gt 0) {
            $firstHeading + 1
        } else {
            0
        }
    }

    # 找到所有旧目录位置
    $tocRanges = Find-TOCRanges -Lines $lines

    # 重建内容
    $newLines = @()
    $tocInserted = $false
    $prevLine = ""

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        $trimmed = $line.Trim()

        # 检查是否在旧目录区域内
        $inTOC = $false
        foreach ($range in $tocRanges) {
            if ($i -ge $range.Start -and $i -lt $range.End) {
                $inTOC = $true
                break
            }
        }

        # 跳过旧目录区域
        if ($inTOC) {
            continue
        }

        # 清理连续的分隔线
        if ($trimmed -eq '---' -and $prevLine.Trim() -eq '---') {
            $prevLine = $line
            continue
        }

        # 在插入位置添加新目录
        if (-not $tocInserted -and $i -eq $insertPos) {
            $newLines += $newTOC.TrimEnd()
            $newLines += ""
            $newLines += "---"
            $newLines += ""
            $tocInserted = $true
        }

        $newLines += $line
        $prevLine = $line
    }

    # 如果还没插入，在末尾添加
    if (-not $tocInserted) {
        $newLines += ""
        $newLines += $newTOC.TrimEnd()
        $newLines += ""
        $newLines += "---"
        $newLines += ""
    }

    return $newLines -join "`n"
}

# 检查文件是否需要修复
function Check-File {
    param([string]$Content)

    $issues = @()

    # 检查目录数量
    $tocCount = ([regex]::Matches($Content, '##\s+📋\s+目录|##\s+目录|#\s+目录')).Count

    # 提取标题
    $headings = Get-Headings -Content $Content
    $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

    # 如果标题少于2个，不需要目录
    if ($filteredHeadings.Count -lt 2) {
        return @{
            NeedsFix = $false
            Issues = @()
        }
    }

    # 检查问题
    if ($tocCount -eq 0) {
        $issues += "缺少目录"
    } elseif ($tocCount -gt 1) {
        $issues += "有 $tocCount 个目录（应只有1个）"
    }

    # 检查重复分隔线
    if ($Content -match '---\s*\n\s*---\s*\n\s*---') {
        $issues += "有重复分隔线"
    }

    # 检查目录位置（应该在元数据后）
    if ($tocCount -gt 0) {
        $lines = $Content -split "`n"
        $metadataEnd = -1
        $tocPos = -1

        for ($i = 0; $i -lt $lines.Count; $i++) {
            $trimmed = $lines[$i].Trim()
            if ($trimmed -eq '---' -and $i -gt 0 -and $metadataEnd -eq -1) {
                $metadataEnd = $i
            }
            if ($trimmed -match '^##+\s+.*目录' -and $tocPos -eq -1) {
                $tocPos = $i
            }
        }

        # 如果目录在元数据之前，需要修复
        if ($metadataEnd -gt 0 -and $tocPos -gt 0 -and $tocPos -lt $metadataEnd) {
            $issues += "目录位置不正确（应在元数据后）"
        }
    }

    return @{
        NeedsFix = $issues.Count -gt 0
        Issues = $issues
    }
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse -File
$totalFiles = $mdFiles.Count
$processedFiles = 0
$fixedFiles = 0
$skippedFiles = 0
$allIssues = @()

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $checkResult = Check-File -Content $content

        if ($checkResult.NeedsFix) {
            $relativePath = $file.FullName.Replace((Get-Location).Path + "\", "")

            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Cyan
            foreach ($issue in $checkResult.Issues) {
                Write-Host "  问题: $issue" -ForegroundColor Yellow
            }

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

            $allIssues += @{
                File = $relativePath
                Issues = $checkResult.Issues
            }
        } else {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 正常" -ForegroundColor Gray
            }
            $skippedFiles++
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
Write-Host "  发现问题: $($allIssues.Count)" -ForegroundColor Yellow
Write-Host "  已修复: $fixedFiles" -ForegroundColor Green
Write-Host "  跳过: $skippedFiles" -ForegroundColor Gray

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际修复" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

if ($allIssues.Count -gt 0 -and $Verbose) {
    Write-Host ""
    Write-Host "问题文件列表:" -ForegroundColor Yellow
    $allIssues | ForEach-Object {
        Write-Host "  - $($_.File)" -ForegroundColor Gray
        foreach ($issue in $_.Issues) {
            Write-Host "    • $issue" -ForegroundColor DarkGray
        }
    }
}

Write-Host ""
