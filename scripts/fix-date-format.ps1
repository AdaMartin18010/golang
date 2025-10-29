# 修复日期格式脚本
# 统一 "2025年XX月XX日" -> "2025-10-29"

param(
    [string]$Path = "docs",
    [switch]$Verbose
)

$stats = @{ Fixed = 0; Scanned = 0 }

Write-Host "🔧 开始修复日期格式..." -ForegroundColor Cyan

$files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File

foreach ($file in $files) {
    $stats.Scanned++
    $content = Get-Content $file.FullName -Raw -Encoding UTF8
    $modified = $false
    
    # 修复: "最后更新**: 2025年10月22日" -> "最后更新**: 2025-10-29"
    if ($content -match '\*\*最后更新\*\*:\s*2025年\d+月\d+日') {
        $content = $content -replace '\*\*最后更新\*\*:\s*2025年\d+月\d+日', '**最后更新**: 2025-10-29'
        $modified = $true
    }
    
    # 修复: "更新**: 2025-10-22" -> "更新**: 2025-10-29"
    if ($content -match '\*\*更新\*\*:\s*2025-10-\d{2}') {
        $content = $content -replace '\*\*更新\*\*:\s*2025-10-\d{2}', '**更新**: 2025-10-29'
        $modified = $true
    }
    
    if ($modified) {
        Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
        $stats.Fixed++
        if ($Verbose) {
            Write-Host "  ✓ $($file.Name)" -ForegroundColor Green
        }
    }
}

Write-Host "`n📊 统计:" -ForegroundColor Cyan
Write-Host "  扫描: $($stats.Scanned)" -ForegroundColor Gray
Write-Host "  修复: $($stats.Fixed)" -ForegroundColor Green

