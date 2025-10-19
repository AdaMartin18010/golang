# Model目录文档序号修正PowerShell脚本
# 用于批量修正所有markdown文档的序号错误

Write-Host "🚀 开始修正Model目录文档序号..." -ForegroundColor Green

# 定义需要修正的模式
$patterns = @{
    "^# 1 1 1 1 1 1 1" = "# "
    "^## 1 1 1 1 1 1 1" = "## 1. "
    "^### 1 1 1 1 1 1 1" = "### 1.1 "
    "^#### 1 1 1 1 1 1 1" = "#### 1.1.1 "
    "^## 9 9 9 9 9 9 9" = "## 2. "
    "^### 9 9 9 9 9 9 9" = "### 2.1 "
    "^## 13 13 13 13 13 13 13" = "## 3. "
    "^### 13 13 13 13 13 13 13" = "### 3.1 "
    "^## 14 14 14 14 14 14 14" = "## 4. "
    "^### 14 14 14 14 14 14 14" = "### 4.1 "
    "^## 15 15 15 15 15 15 15" = "## 5. "
    "^### 15 15 15 15 15 15 15" = "### 5.1 "
    "^## 7 7 7 7 7 7 7" = "## 6. "
    "^### 7 7 7 7 7 7 7" = "### 6.1 "
    "^## 8 8 8 8 8 8 8" = "## 7. "
    "^### 8 8 8 8 8 8 8" = "### 7.1 "
    "^## 11 11 11 11 11 11 11" = "## 8. "
    "^### 11 11 11 11 11 11 11" = "### 8.1 "
    "^## 12 12 12 12 12 12 12" = "## 9. "
    "^### 12 12 12 12 12 12 12" = "### 9.1 "
}

# 统计变量
$totalFiles = 0
$fixedFiles = 0
$skippedFiles = 0

# 获取所有markdown文件
$markdownFiles = Get-ChildItem -Path "." -Recurse -Filter "*.md" -File

foreach ($file in $markdownFiles) {
    $totalFiles++
    Write-Host "📝 处理文件: $($file.FullName)" -ForegroundColor Yellow
    
    # 读取文件内容
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content
    
    # 检查是否需要修正
    $needsFix = $false
    foreach ($pattern in $patterns.Keys) {
        if ($content -match $pattern) {
            $needsFix = $true
            break
        }
    }
    
    if ($needsFix) {
        Write-Host "  🔧 发现序号错误，开始修正..." -ForegroundColor Red
        
        # 创建备份
        $backupPath = "$($file.FullName).bak"
        Copy-Item -Path $file.FullName -Destination $backupPath
        
        # 应用所有修正规则
        foreach ($pattern in $patterns.Keys) {
            $replacement = $patterns[$pattern]
            $content = $content -replace $pattern, $replacement
        }
        
        # 写回文件
        Set-Content -Path $file.FullName -Value $content -Encoding UTF8
        
        # 验证修正结果
        $remainingErrors = $content -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13"
        if ($remainingErrors) {
            Write-Host "  ⚠️  仍有序号错误，请手动检查" -ForegroundColor Yellow
        } else {
            Write-Host "  ✅ 序号修正完成" -ForegroundColor Green
            $fixedFiles++
            # 删除备份文件
            Remove-Item -Path $backupPath
        }
    } else {
        Write-Host "  ✅ 序号格式正确，跳过" -ForegroundColor Green
        $skippedFiles++
    }
    
    Write-Host ""
}

Write-Host "📊 修正完成统计:" -ForegroundColor Cyan
Write-Host "  总文件数: $totalFiles" -ForegroundColor White
Write-Host "  修正文件数: $fixedFiles" -ForegroundColor Green
Write-Host "  跳过文件数: $skippedFiles" -ForegroundColor Blue
Write-Host ""

# 检查是否还有序号错误
Write-Host "🔍 检查剩余序号错误..." -ForegroundColor Cyan
$remainingErrorFiles = Get-ChildItem -Path "." -Recurse -Filter "*.md" -File | Where-Object {
    $content = Get-Content -Path $_.FullName -Raw
    $content -match "1 1 1 1 1 1 1|9 9 9 9 9 9 9|13 13 13 13 13 13 13"
}

if ($remainingErrorFiles.Count -gt 0) {
    Write-Host "⚠️  仍有 $($remainingErrorFiles.Count) 个文件存在序号错误，需要手动处理:" -ForegroundColor Yellow
    $remainingErrorFiles | ForEach-Object { Write-Host "  $($_.FullName)" -ForegroundColor Red }
} else {
    Write-Host "🎉 所有文档序号修正完成！" -ForegroundColor Green
}

Write-Host ""
Write-Host "✨ 序号修正脚本执行完成" -ForegroundColor Green
