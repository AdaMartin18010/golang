# 修复多余空行脚本
# 功能: 修复Markdown文件中的多余空行（MD012规则）
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$DryRun = $false
)

Write-Host "🔧 修复多余空行脚本" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

function Fix-BlankLines {
    param([string]$Content)

    $lines = $Content -split "`n"
    $newLines = @()
    $prevWasBlank = $false

    foreach ($line in $lines) {
        $trimmed = $line.Trim()
        $isBlank = $trimmed -eq ''

        # 如果当前行是空行
        if ($isBlank) {
            # 如果前一行也是空行，跳过（只保留一个空行）
            if ($prevWasBlank) {
                continue
            }
            $prevWasBlank = $true
        } else {
            $prevWasBlank = $false
        }

        $newLines += $line
    }

    return $newLines -join "`n"
}

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
        $newContent = Fix-BlankLines -Content $content

        if ($newContent -ne $content) {
            Write-Host "[$processedFiles/$totalFiles] 🔧 $($file.Name)" -ForegroundColor Yellow

            if (-not $DryRun) {
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

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 修复总结" -ForegroundColor Cyan
Write-Host "  处理文件: $processedFiles/$totalFiles" -ForegroundColor White
Write-Host "  已修复: $fixedFiles" -ForegroundColor Green

if ($DryRun) {
    Write-Host "  ⚠ 这是预览模式，未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "  ✅ 文件已更新" -ForegroundColor Green
}

Write-Host ""
