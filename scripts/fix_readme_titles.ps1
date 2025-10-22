# 修复README文件缺少标题的问题
# 用途: 为缺少一级标题的README.md文件添加合适的标题

param(
    [string]$TargetDir = "docs-new",
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

# 统计
$total = 0
$fixed = 0

Write-Host "=== README标题修复工具 ===" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "⚠️ 演练模式 - 不会实际修改文件" -ForegroundColor Yellow
    Write-Host ""
}

# 需要修复的README文件列表
$targetFiles = @(
    "05-微服务架构\README.md",
    "07-性能优化\README.md"
)

foreach ($relPath in $targetFiles) {
    $filePath = Join-Path $TargetDir $relPath
    
    if (-not (Test-Path $filePath)) {
        Write-Host "  ⚠️ 文件不存在: $relPath" -ForegroundColor Yellow
        continue
    }
    
    $total++
    
    # 读取文件内容
    $content = Get-Content -Path $filePath -Raw
    
    # 检查是否已有一级标题
    if ($content -match "^#\s+") {
        Write-Host "  跳过: $relPath (已有标题)" -ForegroundColor Gray
        continue
    }
    
    # 根据目录名生成标题
    $dirName = Split-Path -Leaf (Split-Path -Parent $filePath)
    $title = "# $dirName"
    
    # 添加标题到文件开头
    $newContent = $title + "`n`n" + $content
    
    if (-not $DryRun) {
        Set-Content -Path $filePath -Value $newContent -NoNewline
        Write-Host "  ✅ 已修复: $relPath" -ForegroundColor Green
        Write-Host "     标题: $title" -ForegroundColor Gray
    } else {
        Write-Host "  [演练] 将修复: $relPath" -ForegroundColor Yellow
        Write-Host "     标题: $title" -ForegroundColor Gray
    }
    
    $fixed++
}

Write-Host ""
Write-Host "=== 完成 ===" -ForegroundColor Cyan
Write-Host "  检查文件: $total" -ForegroundColor White
Write-Host "  已修复: $fixed" -ForegroundColor Green

