# 检查缺少目录的文件脚本
# 功能: 检查哪些Markdown文件缺少目录
# 日期: 2025-11-11

param(
    [string]$DocsPath = "docs",
    [switch]$Verbose = $false
)

Write-Host "📋 检查缺少目录的文件" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

$files = Get-ChildItem -Path $DocsPath -Filter "*.md" -Recurse -File
$totalFiles = $files.Count
$missingTOC = @()
$hasTOC = 0

Write-Host "📁 找到 $totalFiles 个Markdown文件" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $files) {
    try {
        $content = Get-Content $file.FullName -Raw -Encoding UTF8 -ErrorAction SilentlyContinue

        if (-not $content) {
            continue
        }

        # 检查是否有标题
        if ($content -notmatch '^##|^#') {
            continue
        }

        # 检查是否有目录
        $hasTOCPattern = $content -match '##\s+📋\s+目录|##\s+目录|#\s+目录'

        if ($hasTOCPattern) {
            $hasTOC++
            if ($Verbose) {
                Write-Host "✓ $($file.Name)" -ForegroundColor Gray
            }
        } else {
            $relativePath = $file.FullName.Replace((Get-Location).Path + "\", "")
            $missingTOC += $relativePath

            if (-not $Verbose) {
                Write-Host "✗ $relativePath" -ForegroundColor Yellow
            }
        }
    } catch {
        Write-Host "❌ 错误处理 $($file.Name): $_" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "📊 统计结果" -ForegroundColor Cyan
Write-Host "  总文件数: $totalFiles" -ForegroundColor White
Write-Host "  有目录: $hasTOC" -ForegroundColor Green
Write-Host "  缺少目录: $($missingTOC.Count)" -ForegroundColor Yellow

if ($missingTOC.Count -gt 0) {
    $coverage = [math]::Round($hasTOC / $totalFiles * 100, 1)
    Write-Host "  目录覆盖率: $coverage%" -ForegroundColor $(if ($coverage -ge 95) { "Green" } elseif ($coverage -ge 80) { "Yellow" } else { "Red" })

    Write-Host ""
    Write-Host "缺少目录的文件列表:" -ForegroundColor Yellow
    $missingTOC | ForEach-Object { Write-Host "  - $_" -ForegroundColor Gray }
} else {
    Write-Host ""
    Write-Host "✅ 所有文件都有目录！" -ForegroundColor Green
}

Write-Host ""
