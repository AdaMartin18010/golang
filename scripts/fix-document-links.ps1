# 文档内部链接修复脚本
# 功能: 修复Markdown文档中的内部链接失效问题
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false,
    [switch]$Verbose = $false
)

Write-Host "🔧 文档内部链接修复脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# GitHub Markdown anchor 生成规则
function Get-GitHubAnchor {
    param([string]$Title)

    $anchor = $Title

    # 1. 转换为小写
    $anchor = $anchor.ToLower()

    # 2. 移除emoji（保留中文字符）
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
    param([string]$FilePath)

    $headings = @()
    $content = Get-Content $FilePath -Raw -Encoding UTF8

    # 匹配所有标题 (# 标题)
    $matches = [regex]::Matches($content, '^(#{1,6})\s+(.+)$', [System.Text.RegularExpressions.RegexOptions]::Multiline)

    foreach ($match in $matches) {
        $level = $match.Groups[1].Value.Length
        $title = $match.Groups[2].Value.Trim()
        $anchor = Get-GitHubAnchor $title

        $headings += @{
            Level = $level
            Title = $title
            Anchor = $anchor
            Line = $match.Groups[0].Value
        }
    }

    return $headings
}

# 修复文档中的链接
function Fix-DocumentLinks {
    param(
        [string]$FilePath,
        [array]$Headings
    )

    $content = Get-Content $FilePath -Raw -Encoding UTF8
    $originalContent = $content
    $fixCount = 0

    # 匹配所有内部链接 [文本](#anchor)
    $linkPattern = '\[([^\]]+)\]\(#([^\)]+)\)'
    $matches = [regex]::Matches($content, $linkPattern)

    foreach ($match in $matches) {
        $linkText = $match.Groups[1].Value
        $oldAnchor = $match.Groups[2].Value

        # 查找匹配的标题
        $matchedHeading = $Headings | Where-Object {
            $_.Title -eq $linkText -or
            $_.Anchor -eq $oldAnchor -or
            (Get-GitHubAnchor $_.Title) -eq $oldAnchor
        } | Select-Object -First 1

        if ($matchedHeading) {
            $newAnchor = $matchedHeading.Anchor
            if ($newAnchor -ne $oldAnchor) {
                $newLink = "[$linkText](#$newAnchor)"
                $content = $content -replace [regex]::Escape($match.Value), $newLink
                $fixCount++

                if ($Verbose) {
                    Write-Host "  ✓ 修复: [$linkText](#$oldAnchor) -> [$linkText](#$newAnchor)" -ForegroundColor Green
                }
            }
        } else {
            if ($Verbose) {
                Write-Host "  ⚠ 未找到匹配: [$linkText](#$oldAnchor)" -ForegroundColor Yellow
            }
        }
    }

    return @{
        Content = $content
        FixCount = $fixCount
        Changed = $content -ne $originalContent
    }
}

# 主处理逻辑
$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse
$totalFiles = $mdFiles.Count
$totalFixes = 0
$processedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++
    Write-Host "[$processedFiles/$totalFiles] 处理: $($file.FullName)" -ForegroundColor Cyan

    try {
        # 提取标题
        $headings = Get-Headings -FilePath $file.FullName

        if ($headings.Count -eq 0) {
            Write-Host "  ⚠ 未找到标题，跳过" -ForegroundColor Yellow
            continue
        }

        # 修复链接
        $result = Fix-DocumentLinks -FilePath $file.FullName -Headings $headings

        if ($result.FixCount -gt 0) {
            Write-Host "  ✓ 修复了 $($result.FixCount) 个链接" -ForegroundColor Green
            $totalFixes += $result.FixCount

            if (-not $DryRun -and $result.Changed) {
                # 保存文件
                [System.IO.File]::WriteAllText($file.FullName, $result.Content, [System.Text.Encoding]::UTF8)
            }
        } else {
            Write-Host "  ✓ 无需修复" -ForegroundColor Gray
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }

    Write-Host ""
}

# 总结
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 修复总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  修复链接: $totalFixes" -ForegroundColor White

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:$false 来实际修复" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""
