# 高级序号修正PowerShell脚本
# 处理复杂的序号错误模式

Write-Host "🔧 开始高级序号修正..." -ForegroundColor Green

# 更精确的修正模式
$advancedPatterns = @{
    # 处理标题中的重复数字
    "^# (\d+) \1 \1 \1 \1 \1 \1" = "# "
    "^## (\d+) \1 \1 \1 \1 \1 \1" = "## $1. "
    "^### (\d+) \1 \1 \1 \1 \1 \1" = "### $1.1 "
    "^#### (\d+) \1 \1 \1 \1 \1 \1" = "#### $1.1.1 "
    
    # 处理TOC中的错误链接
    "- \[(\d+) \1 \1 \1 \1 \1 \1" = "- [$1"
    "- \[(\d+) \1 \1 \1 \1 \1 \1 ([^\]]+)\]" = "- [$1 $2]"
    
    # 处理内容中的重复数字
    "## (\d+) \1 \1 \1 \1 \1 \1" = "## $1. "
    "### (\d+) \1 \1 \1 \1 \1 \1" = "### $1.1 "
    "#### (\d+) \1 \1 \1 \1 \1 \1" = "#### $1.1.1 "
}

# 获取仍有错误的文件
$errorFiles = Get-ChildItem -Path "." -Recurse -Filter "*.md" -File | Where-Object {
    $content = Get-Content -Path $_.FullName -Raw
    $content -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13|14 14 14 14 14 14 14|15 15 15 15 15 15 15"
}

Write-Host "发现 $($errorFiles.Count) 个文件需要高级修正" -ForegroundColor Yellow

$fixedCount = 0
$skippedCount = 0

foreach ($file in $errorFiles) {
    Write-Host "📝 处理文件: $($file.FullName)" -ForegroundColor Cyan
    
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content
    
    # 应用高级修正模式
    foreach ($pattern in $advancedPatterns.Keys) {
        $replacement = $advancedPatterns[$pattern]
        $content = $content -replace $pattern, $replacement
    }
    
    # 检查是否还有错误
    $stillHasErrors = $content -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13|14 14 14 14 14 14 14|15 15 15 15 15 15 15"
    
    if ($stillHasErrors) {
        Write-Host "  ⚠️  仍有序号错误，需要手动处理" -ForegroundColor Red
        
        # 显示具体的错误行
        $lines = $content -split "`n"
        for ($i = 0; $i -lt $lines.Count; $i++) {
            if ($lines[$i] -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13|14 14 14 14 14 14 14|15 15 15 15 15 15 15") {
                Write-Host "    行 $($i + 1): $($lines[$i].Trim())" -ForegroundColor Red
            }
        }
        $skippedCount++
    } else {
        # 写回修正后的内容
        Set-Content -Path $file.FullName -Value $content -Encoding UTF8
        Write-Host "  ✅ 序号修正完成" -ForegroundColor Green
        $fixedCount++
    }
    
    Write-Host ""
}

Write-Host "📊 高级修正完成统计:" -ForegroundColor Cyan
Write-Host "  修正文件数: $fixedCount" -ForegroundColor Green
Write-Host "  跳过文件数: $skippedCount" -ForegroundColor Yellow
Write-Host ""

# 最终检查
$finalErrorFiles = Get-ChildItem -Path "." -Recurse -Filter "*.md" -File | Where-Object {
    $content = Get-Content -Path $_.FullName -Raw
    $content -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13|14 14 14 14 14 14 14|15 15 15 15 15 15 15"
}

if ($finalErrorFiles.Count -eq 0) {
    Write-Host "🎉 所有文档序号修正完成！" -ForegroundColor Green
} else {
    Write-Host "⚠️  仍有 $($finalErrorFiles.Count) 个文件需要手动处理" -ForegroundColor Yellow
    Write-Host "这些文件可能需要更复杂的修正逻辑或内容重构" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "✨ 高级序号修正脚本执行完成" -ForegroundColor Green
