# 全面目录检查脚本
# 功能: 检查所有Markdown文件的目录质量和一致性
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$Fix = $false
)

Write-Host "📋 全面目录检查脚本" -ForegroundColor Cyan
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

# 检查目录
function Check-TOC {
    param([string]$Content, [string]$FilePath)

    $issues = @()
    $lines = $Content -split "`n"

    # 检查目录数量
    $tocCount = ([regex]::Matches($Content, '##\s+📋\s+目录|##\s+目录|#\s+目录')).Count

    if ($tocCount -eq 0) {
        $issues += "缺少目录"
    } elseif ($tocCount -gt 1) {
        $issues += "有 $tocCount 个目录（应只有1个）"
    }

    # 提取标题
    $headings = Get-Headings -Content $Content
    $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

    if ($filteredHeadings.Count -lt 2) {
        return @{
            HasIssues = $false
            Issues = @()
            NeedsTOC = $false
        }
    }

    # 检查目录项是否完整
    if ($tocCount -gt 0) {
        # 提取目录中的链接
        $tocLinks = [regex]::Matches($Content, '\[([^\]]+)\]\(#([^\)]+)\)')
        $tocAnchors = $tocLinks | ForEach-Object { $_.Groups[2].Value }

        # 检查每个标题是否在目录中
        foreach ($heading in $filteredHeadings) {
            if ($heading.Anchor -notin $tocAnchors) {
                $issues += "标题 '$($heading.Title)' 不在目录中"
            }
        }
    }

    return @{
        HasIssues = $issues.Count -gt 0
        Issues = $issues
        NeedsTOC = $filteredHeadings.Count -ge 2
        Headings = $filteredHeadings
    }
}

# 生成目录
function Generate-TOC {
    param([array]$Headings)
    if ($Headings.Count -eq 0) { return "" }

    $toc = "## 📋 目录`n`n"
    foreach ($heading in $Headings) {
        $level = $heading.Level
        if ($level -gt 2) { continue }
        $indent = if ($level -gt 1) { "  " * ($level - 1) } else { "" }
        $toc += "$indent- [$($heading.Title)](#$($heading.Anchor))`n"
    }
    return $toc
}

# 修复目录
function Fix-TOC {
    param([string]$Content, [array]$Headings)

    $lines = $Content -split "`n"
    $newLines = @()
    $tocInserted = $false

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

    # 生成新目录
    $newTOC = Generate-TOC -Headings $Headings

    # 找到所有目录位置
    $tocStarts = @()
    for ($i = 0; $i -lt $lines.Count; $i++) {
        if ($lines[$i].Trim() -match '^##+\s+.*目录') {
            $tocStarts += $i
        }
    }

    # 重建内容
    for ($i = 0; $i -lt $lines.Count; $i++) {
        # 跳过所有旧目录区域
        $inTOC = $false
        foreach ($tocStart in $tocStarts) {
            if ($i -ge $tocStart) {
                # 检查是否在目录区域内
                $j = $i
                while ($j -lt $lines.Count) {
                    $trimmed = $lines[$j].Trim()
                    if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                        break
                    }
                    if ($j -eq $i) {
                        $inTOC = $true
                    }
                    $j++
                }
                if ($inTOC) {
                    # 跳过目录项
                    if ($trimmed -match '^\s*-\s+\[.*\]\(#.*\)') {
                        continue
                    }
                    # 跳过目录标题和分隔线
                    if ($trimmed -match '^##+\s+.*目录' -or $trimmed -eq '---' -or $trimmed -eq '') {
                        continue
                    }
                    # 如果遇到下一个标题，停止跳过
                    if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                        $inTOC = $false
                    }
                }
            }
        }

        if ($inTOC) {
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

        $newLines += $lines[$i]
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
$filesWithIssues = 0
$fixedFiles = 0
$allIssues = @()

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $checkResult = Check-TOC -Content $content -FilePath $file.FullName

        if ($checkResult.NeedsTOC -and $checkResult.HasIssues) {
            $filesWithIssues++
            $relativePath = $file.FullName.Replace((Get-Location).Path + "\", "")

            Write-Host "[$processedFiles/$totalFiles] ⚠️ $($file.Name)" -ForegroundColor Yellow
            foreach ($issue in $checkResult.Issues) {
                Write-Host "  问题: $issue" -ForegroundColor Red
            }

            $allIssues += @{
                File = $relativePath
                Issues = $checkResult.Issues
            }

            if ($Fix) {
                $newContent = Fix-TOC -Content $content -Headings $checkResult.Headings
                if ($newContent -ne $content) {
                    [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                    Write-Host "  ✅ 已修复" -ForegroundColor Green
                    $fixedFiles++
                }
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
Write-Host "  发现问题: $filesWithIssues" -ForegroundColor Yellow
if ($Fix) {
    Write-Host "  已修复: $fixedFiles" -ForegroundColor Green
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
} else {
    Write-Host "  ⚠ 这是检查模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -Fix 来实际修复问题" -ForegroundColor Yellow
}

if ($allIssues.Count -gt 0) {
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
