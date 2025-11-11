# 补充缺失目录脚本
# 功能: 为缺少目录的文档自动生成目录
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "📋 补充缺失目录脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# GitHub Markdown anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Title)

    $anchor = $Title

    # 1. 转换为小写
    $anchor = $anchor.ToLower()

    # 2. 移除emoji和特殊字符（保留中文字符）
    $anchor = $anchor -replace '[^\w\s\u4e00-\u9fa5-]', ''

    # 3. 替换空格为连字符
    $anchor = $anchor -replace '\s+', '-'

    # 4. 移除多余的连字符
    $anchor = $anchor -replace '-+', '-'

    # 5. 移除首尾连字符
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

        # 匹配标题 (# 标题)
        if ($trimmed -match '^(#{1,6})\s+(.+)$') {
            $level = $matches[1].Length
            $title = $matches[2].Trim()

            # 跳过目录标题本身
            if ($title -match '目录') {
                continue
            }

            $anchor = Get-GitHubAnchor $title

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

# 检测是否已有目录
function Has-TOC {
    param([string]$Content)

    $lines = $Content -split "`n"
    $foundTocTitle = $false

    foreach ($line in $lines) {
        $trimmed = $line.Trim()

        # 检测目录标题
        if ($trimmed -match '^##+\s+.*目录') {
            $foundTocTitle = $true
            continue
        }

        # 如果找到目录标题，检查是否有目录项
        if ($foundTocTitle) {
            if ($trimmed -match '^\s*-\s+\[.*\]\(#.*\)') {
                return $true
            }
            # 如果遇到下一个标题，停止检查
            if ($trimmed -match '^##+\s+') {
                break
            }
        }
    }

    return $false
}

# 插入目录
function Insert-TOC {
    param(
        [string]$Content,
        [string]$TOC
    )

    $lines = $Content -split "`n"
    $newLines = @()
    $inserted = $false

    # 查找插入位置（元数据后，第一个内容标题前）
    $metadataEnd = -1
    $firstHeading = -1

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i].Trim()

        # 检测分隔线（元数据结束）
        if ($line -eq '---' -and $metadataEnd -eq -1) {
            $metadataEnd = $i
            continue
        }

        # 检测第一个内容标题
        if ($metadataEnd -gt 0 -and $firstHeading -eq -1) {
            if ($line -match '^##+\s+' -and $line -notmatch '目录') {
                $firstHeading = $i
                break
            }
        }
    }

    # 确定插入位置
    $insertPos = if ($metadataEnd -gt 0) { $metadataEnd + 1 } else { 0 }

    # 重建内容
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

    return $newLines -join "`n"
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$processedFiles = 0
$addedTocFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8

        # 检查是否已有目录
        if (Has-TOC -Content $content) {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] ✓ $($file.Name) - 已有目录" -ForegroundColor Gray
            }
            continue
        }

        # 提取标题
        $headings = Get-Headings -Content $content

        # 过滤：只保留一级和二级标题（避免目录过长）
        $filteredHeadings = $headings | Where-Object { $_.Level -le 2 }

        if ($filteredHeadings.Count -lt 2) {
            if ($Verbose) {
                Write-Host "[$processedFiles/$totalFiles] - $($file.Name) - 标题太少，跳过" -ForegroundColor Gray
            }
            continue
        }

        Write-Host "[$processedFiles/$totalFiles] 📋 $($file.Name)" -ForegroundColor Cyan
        Write-Host "  标题数: $($filteredHeadings.Count)" -ForegroundColor Yellow

        # 生成目录
        $toc = Generate-TOC -Headings $filteredHeadings

        if (-not $DryRun) {
            $newContent = Insert-TOC -Content $content -TOC $toc
            [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
            Write-Host "  ✅ 已添加目录" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
        }

        $addedTocFiles++
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

# 总结
Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 补充总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已添加目录: $addedTocFiles" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际添加目录" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""
