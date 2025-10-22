# 批量补充文档元信息
# 用途: 为docs-new目录下的所有md文件添加统一的元信息

param(
    [string]$TargetDir = "docs-new",
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

# 元信息模板
$metadataTemplate = @"

---

**维护者**: Documentation Team  
**创建日期**: 2025-10-22  
**最后更新**: 2025-10-22  
**文档状态**: ✅ 完成
"@

# 统计
$total = 0
$updated = 0
$skipped = 0

Write-Host "=== 文档元信息补充工具 ===" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "⚠️ 演练模式 - 不会实际修改文件" -ForegroundColor Yellow
    Write-Host ""
}

# 获取所有md文件
$files = Get-ChildItem -Path $TargetDir -Filter "*.md" -Recurse

foreach ($file in $files) {
    $total++
    
    # 读取文件内容
    $content = Get-Content -Path $file.FullName -Raw
    
    # 检查是否已有元信息
    if ($content -match "维护者.*Documentation Team") {
        Write-Host "  跳过: $($file.Name) (已有元信息)" -ForegroundColor Gray
        $skipped++
        continue
    }
    
    # 添加元信息
    $newContent = $content.TrimEnd() + $metadataTemplate
    
    if (-not $DryRun) {
        Set-Content -Path $file.FullName -Value $newContent -NoNewline
        Write-Host "  ✅ 已更新: $($file.Name)" -ForegroundColor Green
    } else {
        Write-Host "  [演练] 将更新: $($file.Name)" -ForegroundColor Yellow
    }
    
    $updated++
}

Write-Host ""
Write-Host "=== 完成 ===" -ForegroundColor Cyan
Write-Host "  总文件数: $total" -ForegroundColor White
Write-Host "  已更新: $updated" -ForegroundColor Green
Write-Host "  已跳过: $skipped" -ForegroundColor Gray

