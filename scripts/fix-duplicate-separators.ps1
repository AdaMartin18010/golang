# 修复重复分隔线脚本
# 功能: 修复Markdown文件中的重复分隔线
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false
)

Write-Host "🔧 修复重复分隔线脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

$mdFiles = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse -File
$totalFiles = $mdFiles.Count
$processedFiles = 0
$fixedFiles = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $mdFiles) {
    $processedFiles++

    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        $lines = $content -split "`n"
        $newLines = @()
        $prevLine = ""
        $modified = $false

        foreach ($line in $lines) {
            $trimmed = $line.Trim()

            # 检查连续的分隔线
            if ($trimmed -eq '---' -and $prevLine.Trim() -eq '---') {
                $modified = $true
                continue
            }

            $newLines += $line
            $prevLine = $line
        }

        if ($modified) {
            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Cyan

            if (-not $DryRun) {
                $newContent = $newLines -join "`n"
                [System.IO.File]::WriteAllText($file.FullName, $newContent, [System.Text.Encoding]::UTF8)
                Write-Host "  ✅ 已修复" -ForegroundColor Green
                $fixedFiles++
            } else {
                Write-Host "  ⚠ 预览模式（未实际修改）" -ForegroundColor Yellow
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
Write-Host "  已修复: $fixedFiles" -ForegroundColor Green

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
    Write-Host "  使用 -DryRun:`$false 来实际修复" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""
