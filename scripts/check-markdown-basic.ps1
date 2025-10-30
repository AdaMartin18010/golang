# Markdown 基础格式检查脚本（无需外部依赖）
# 检查常见的 Markdown 格式问题

param (
    [string]$Path = "docs",
    [switch]$Fix = $false
)

Write-Host "🔍 扫描 $Path 目录中的所有 .md 文件..." -ForegroundColor Cyan
Write-Host ""

$markdownFiles = Get-ChildItem -Path $Path -Filter "*.md" -Recurse
Write-Host "📝 找到 $($markdownFiles.Count) 个 Markdown 文件" -ForegroundColor Cyan
Write-Host ""

$issuesFound = 0
$filesWithIssues = 0
$fixedIssues = 0

foreach ($file in $markdownFiles) {
    $content = Get-Content $file.FullName -Raw
    $originalContent = $content
    $fileIssues = @()
    $lines = Get-Content $file.FullName

    # 检查 1: MD012 - 多个连续空行
    $multipleBlankLines = [regex]::Matches($content, '\n\n\n+')
    if ($multipleBlankLines.Count -gt 0) {
        $fileIssues += "MD012: 发现 $($multipleBlankLines.Count) 处多个连续空行"
        if ($Fix) {
            $content = [regex]::Replace($content, '\n\n\n+', "`n`n")
            $fixedIssues++
        }
    }

    # 检查 2: MD009 - 行尾空格
    $trailingSpaces = 0
    for ($i = 0; $i -lt $lines.Count; $i++) {
        if ($lines[$i] -match '\s+$') {
            $trailingSpaces++
        }
    }
    if ($trailingSpaces -gt 0) {
        $fileIssues += "MD009: 发现 $trailingSpaces 行包含行尾空格"
        if ($Fix) {
            $content = (($content -split "`n") | ForEach-Object { $_.TrimEnd() }) -join "`n"
            $fixedIssues++
        }
    }

    # 检查 3: MD040 - 代码块缺少语言标识
    $codeBlocksWithoutLang = [regex]::Matches($content, '(?m)^```\s*$')
    if ($codeBlocksWithoutLang.Count -gt 0) {
        $fileIssues += "MD040: 发现 $($codeBlocksWithoutLang.Count) 个代码块缺少语言标识"
    }

    # 检查 4: MD042 - 空链接
    $emptyLinks = [regex]::Matches($content, '\[([^\]]+)\]\(#?\s*\)')
    if ($emptyLinks.Count -gt 0) {
        $fileIssues += "MD042: 发现 $($emptyLinks.Count) 个空链接"
    }

    # 检查 5: MD031 - 代码块周围需要空行
    $codeBlocksWithoutBlankLines = [regex]::Matches($content, '(?m)^[^\n`].*\n```')
    if ($codeBlocksWithoutBlankLines.Count -gt 0) {
        $fileIssues += "MD031: 发现 $($codeBlocksWithoutBlankLines.Count) 个代码块前缺少空行"
    }

    # 检查 6: MD047 - 文件应以换行符结尾
    if (-not $content.EndsWith("`n")) {
        $fileIssues += "MD047: 文件未以换行符结尾"
        if ($Fix) {
            $content += "`n"
            $fixedIssues++
        }
    }

    # 检查 7: 重复的版本信息块
    $versionBlockPattern = '(?s)\*\*版本\*\*:\s*v\d+\.\d+\s*\n\*\*更新日期\*\*:\s*\d{4}-\d{2}-\d{2}\s*\n\*\*适用于\*\*:\s*Go\s*\d+\.\d+\.\d+\+?'
    $versionBlocks = [regex]::Matches($content, $versionBlockPattern)
    if ($versionBlocks.Count -gt 1) {
        $fileIssues += "自定义: 发现 $($versionBlocks.Count) 个重复的版本信息块"
        if ($Fix) {
            # 保留第一个，删除其他的
            for ($i = $versionBlocks.Count - 1; $i -ge 1; $i--) {
                $content = $content.Remove($versionBlocks[$i].Index, $versionBlocks[$i].Length)
            }
            $fixedIssues++
        }
    }

    # 如果发现问题，记录并输出
    if ($fileIssues.Count -gt 0) {
        $filesWithIssues++
        $issuesFound += $fileIssues.Count

        Write-Host "📄 $($file.Name)" -ForegroundColor Yellow
        foreach ($issue in $fileIssues) {
            Write-Host "   ⚠️  $issue" -ForegroundColor Yellow
        }
        Write-Host ""
    }

    # 如果启用了修复，保存修改
    if ($Fix -and $content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -NoNewline
        Write-Host "   ✅ 已自动修复部分问题" -ForegroundColor Green
        Write-Host ""
    }
}

Write-Host ""
Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host "📊 检查完成统计" -ForegroundColor Cyan
Write-Host "=" * 80 -ForegroundColor Cyan
Write-Host "总文件数: $($markdownFiles.Count)" -ForegroundColor White
Write-Host "问题文件: $filesWithIssues" -ForegroundColor $(if ($filesWithIssues -eq 0) { "Green" } else { "Yellow" })
Write-Host "问题总数: $issuesFound" -ForegroundColor $(if ($issuesFound -eq 0) { "Green" } else { "Yellow" })
if ($Fix) {
    Write-Host "已修复数: $fixedIssues" -ForegroundColor Green
}
Write-Host "=" * 80 -ForegroundColor Cyan

if ($issuesFound -gt 0 -and -not $Fix) {
    Write-Host ""
    Write-Host "💡 提示: 使用 -Fix 参数自动修复可修复的问题" -ForegroundColor Cyan
    Write-Host "   示例: .\scripts\check-markdown-basic.ps1 -Path docs -Fix" -ForegroundColor White
}

if ($issuesFound -eq 0) {
    Write-Host ""
    Write-Host "✅ 所有文件检查通过！" -ForegroundColor Green
}

exit $(if ($issuesFound -eq 0) { 0 } else { 1 })
