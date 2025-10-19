# 最终验证PowerShell脚本
# 用于验证所有TOC和序号修正结果

Write-Host "🔍 开始最终验证..." -ForegroundColor Green

# 获取所有markdown文件
$markdownFiles = Get-ChildItem -Path "." -Recurse -Filter "*.md" -File

$totalFiles = 0
$validFiles = 0
$errorFiles = 0
$tocFiles = 0
$noTocFiles = 0

$errorDetails = @()

foreach ($file in $markdownFiles) {
    $totalFiles++
    Write-Host "📝 验证文件: $($file.FullName)" -ForegroundColor Yellow
    
    # 读取文件内容
    $content = Get-Content -Path $file.FullName -Raw
    
    # 检查是否包含TOC
    $hasToc = $content -match "<!-- TOC START -->" -and $content -match "<!-- TOC END -->"
    
    if ($hasToc) {
        $tocFiles++
        Write-Host "  📑 发现TOC" -ForegroundColor Cyan
        
        # 检查TOC格式错误
        $hasErrors = $false
        $errors = @()
        
        # 检查序号错误
        if ($content -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13|14 14 14 14 14 14 14|15 15 15 15 15 15 15") {
            $hasErrors = $true
            $errors += "序号错误"
        }
        
        # 检查TOC链接格式错误
        if ($content -match "1\.2\.\d+|13\.\d+|14\.\d+") {
            $hasErrors = $true
            $errors += "TOC链接格式错误"
        }
        
        # 检查标题格式错误
        if ($content -match "^# 1 1 1 1 1 1 1|^## 9 9 9 9 9 9 9|^## 13 13 13 13 13 13 13") {
            $hasErrors = $true
            $errors += "标题格式错误"
        }
        
        if ($hasErrors) {
            $errorFiles++
            $errorDetails += "$($file.FullName): $($errors -join ', ')"
            Write-Host "  ❌ 发现错误: $($errors -join ', ')" -ForegroundColor Red
        } else {
            $validFiles++
            Write-Host "  ✅ TOC格式正确" -ForegroundColor Green
        }
    } else {
        $noTocFiles++
        Write-Host "  ⚪ 无TOC" -ForegroundColor Gray
    }
    
    Write-Host ""
}

Write-Host "📊 最终验证统计:" -ForegroundColor Cyan
Write-Host "  总文件数: $totalFiles" -ForegroundColor White
Write-Host "  有TOC文件数: $tocFiles" -ForegroundColor Blue
Write-Host "  无TOC文件数: $noTocFiles" -ForegroundColor Gray
Write-Host "  格式正确: $validFiles" -ForegroundColor Green
Write-Host "  格式错误: $errorFiles" -ForegroundColor Red
Write-Host ""

if ($errorFiles -gt 0) {
    Write-Host "⚠️  发现 $errorFiles 个文件存在格式错误:" -ForegroundColor Yellow
    $errorDetails | ForEach-Object { Write-Host "  $_" -ForegroundColor Red }
} else {
    Write-Host "🎉 所有文档格式验证通过！" -ForegroundColor Green
}

Write-Host ""
Write-Host "✨ 最终验证脚本执行完成" -ForegroundColor Green
