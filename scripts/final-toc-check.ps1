# 最终全面检查目录脚本
# 功能: 全面递归检查所有Markdown文件的目录问题
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$Fix = $false
)

Write-Host "📋 最终全面检查目录脚本" -ForegroundColor Cyan
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
        if ($level -gt 2) { continue }
        $indent = if ($level -gt 1) { "  " * ($level - 1) } else { "" }
        $toc += "$indent- [$($heading.Title)](#$($heading.Anchor))`n"
    }
    return $toc
}

# 检查文件
function Check-File {
    param([string]$Content, [string]$FilePath)

    $issues = @()
    $headings = Get-Headings -Content $Content
    $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

    if ($filteredHeadings.Count -lt 2) {
        return @{ NeedsFix = $false; Issues = @() }
    }

    $tocCount = ([regex]::Matches($Content, '##\s+📋\s+目录|##\s+目录|#\s+目录')).Count

    if ($tocCount -eq 0) {
        $issues += "缺少目录"
    } elseif ($tocCount -gt 1) {
        $issues += "有 $tocCount 个目录（应只有1个）"
    }

    # 检查目录位置
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

        # 如果目录在元数据之前或距离元数据太远，需要修复
        if ($metadataEnd -gt 0) {
            if ($tocPos -lt $metadataEnd) {
                $issues += "目录位置不正确（在元数据前）"
            } elseif ($tocPos -gt $metadataEnd + 20) {
                $issues += "目录位置不正确（距离元数据太远）"
            }
        }
    }

    return @{
        NeedsFix = $issues.Count -gt 0
        Issues = $issues
        Headings = $filteredHeadings
    }
}

# 修复文件
function Fix-File {
    param([string]$Content, [array]$Headings)

    $lines = $Content -split "`n"
    $newTOC = Generate-TOC -Headings $Headings

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

    # 找到所有旧目录位置
    $tocRanges = @()
    $tocStart = -1
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $trimmed = $lines[$i].Trim()
        if ($trimmed -match '^##+\s+.*目录') {
            if ($tocStart -eq -1) {
                $tocStart = $i
            }
        } elseif ($tocStart -ge 0) {
            if ($trimmed -match '^##+\s+' -and $trimmed -notmatch '目录') {
                $tocRanges += @{ Start = $tocStart; End = $i }
                $tocStart = -1
            } elseif ($i -eq $lines.Count - 1) {
                $tocRanges += @{ Start = $tocStart; End = $i + 1 }
                $tocStart = -1
            }
        }
    }

    # 重建内容
    $newLines = @()
    $tocInserted = $false

    for ($i = 0; $i -lt $lines.Count; $i++) {
        # 跳过旧目录区域
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
$allIssues = @()

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $checkResult = Check-File -Content $content -FilePath $file.FullName

        if ($checkResult.NeedsFix) {
            $relativePath = $file.FullName.Replace((Get-Location).Path + "\", "")

            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Cyan
            foreach ($issue in $checkResult.Issues) {
                Write-Host "  问题: $issue" -ForegroundColor Yellow
            }

            if ($Fix) {
                $newContent = Fix-File -Content $content -Headings $checkResult.Headings
                if ($newContent -ne $content) {
                    [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                    Write-Host "  ✅ 已修复" -ForegroundColor Green
                    $fixedFiles++
                }
            }

            $allIssues += @{
                File = $relativePath
                Issues = $checkResult.Issues
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
Write-Host "  发现问题: $($allIssues.Count)" -ForegroundColor Yellow
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
