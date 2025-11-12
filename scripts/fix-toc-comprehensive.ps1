# 全面修复目录脚本
# 功能: 确保所有Markdown文件有且只有一个目录，目录完整且有序
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false
)

Write-Host "📋 全面修复目录脚本" -ForegroundColor Cyan
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
        if ($lines[$i].Trim() -eq '---' -and $i -gt 0) {
            $metadataEnd = $i
            break
        }
    }

    # 找到第一个标题位置
    $firstHeading = -1
    for ($i = 0; $i -lt $lines.Count; $i++) {
        if ($lines[$i].Trim() -match '^#\s+') {
            $firstHeading = $i
            break
        }
    }

    # 确定插入位置
    $insertPos = if ($metadataEnd -gt 0) { $metadataEnd + 1 } else { if ($firstHeading -gt 0) { $firstHeading + 1 } else { 0 } }

    # 找到所有目录位置
    $tocRanges = @()
    $tocStart = -1
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()
        if ($trimmed -match '^##+\s+.*目录') {
            if ($tocStart -eq -1) {
                $tocStart = $i
            }
        } elseif ($tocStart -ge 0) {
            # 找到目录结束位置（下一个标题或空行后的标题）
            if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                $tocRanges += @{ Start = $tocStart; End = $i }
                $tocStart = -1
            } elseif ($i -eq $lines.Count - 1) {
                # 文件末尾
                $tocRanges += @{ Start = $tocStart; End = $i + 1 }
                $tocStart = -1
            }
        }
    }

    # 清理重复分隔线
    $newLines = @()
    $prevLine = ""
    $tocInserted = $false

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        $trimmed = $line.Trim()

        # 跳过所有旧目录区域
        $inTOC = $false
        foreach ($range in $tocRanges) {
            if ($i -ge $range.Start -and $i -lt $range.End) {
                $inTOC = $true
                break
            }
        }

        if ($inTOC) {
            continue
        }

        # 清理连续的分隔线
        if ($trimmed -eq '---' -and $prevLine.Trim() -eq '---') {
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
        $headings = Get-Headings -Content $content
        $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

        # 如果标题少于2个，跳过
        if ($filteredHeadings.Count -lt 2) {
            continue
        }

        # 检查目录
        $tocCount = ([regex]::Matches($content, '##\s+📋\s+目录|##\s+目录|#\s+目录')).Count
        $hasDuplicateSeparators = $content -match '---\s*\n\s*---\s*\n\s*---'

        $needsFix = $false
        $issueMsg = ""

        if ($tocCount -eq 0) {
            $needsFix = $true
            $issueMsg = "缺少目录"
        } elseif ($tocCount -gt 1) {
            $needsFix = $true
            $issueMsg = "有 $tocCount 个目录"
        } elseif ($hasDuplicateSeparators) {
            $needsFix = $true
            $issueMsg = "有重复分隔线"
        }

        if ($needsFix) {
            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Cyan
            Write-Host "  问题: $issueMsg" -ForegroundColor Yellow

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

Write-Host ""
