# 清理多余分隔线脚本
# 功能: 清理文档中多余的分隔线，保持结构一致
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false
)

Write-Host "🧹 清理多余分隔线脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

function Clean-Separators {
    param([string]$Content)

    $lines = $Content -split "`n"
    $newLines = @()
    $prevLine = ""
    $prevPrevLine = ""
    $prevPrevPrevLine = ""
    $consecutiveBlanks = 0

    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        $trimmed = $line.Trim()

        # 跳过连续的分隔线
        if ($trimmed -eq '---') {
            # 如果前一行或前前一行也是分隔线，跳过
            if ($prevLine.Trim() -eq '---' -or $prevPrevLine.Trim() -eq '---') {
                continue
            }
        }

        # 处理连续的空行（最多保留一个空行）
        if ($trimmed -eq '') {
            $consecutiveBlanks++
            # 如果已经有空行，跳过这个空行
            if ($consecutiveBlanks > 1) {
                continue
            }
        } else {
            $consecutiveBlanks = 0
        }

        $newLines += $line
        $prevPrevPrevLine = $prevPrevLine
        $prevPrevLine = $prevLine
        $prevLine = $line
    }

    return $newLines -join "`n"
}

$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse -File
$totalFiles = $mdFiles.Count
$processedFiles = 0
$cleanedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $newContent = Clean-Separators -Content $content

        if ($newContent -ne $content) {
            Write-Host "[$processedFiles/$totalFiles] 🧹 $($file.Name)" -ForegroundColor Yellow

            if (-not $DryRun) {
                [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已清理" -ForegroundColor Green
                $cleanedFiles++
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
            }
        }
    } catch {
        Write-Host "  ❌ 错误: $_" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 清理总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已清理: $cleanedFiles" -ForegroundColor Green

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""
