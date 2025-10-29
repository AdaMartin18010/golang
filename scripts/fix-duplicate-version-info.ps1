# 修复文档中重复的版本信息
# 处理模式：
# 1. 连续重复的版本信息块
# 2. 多余的分隔符

$ErrorActionPreference = "Stop"
$filesFixed = 0
$totalIssues = 0

Write-Host "🔍 扫描docs目录中的所有.md文件..." -ForegroundColor Cyan

# 获取所有Markdown文件
$mdFiles = Get-ChildItem -Path "docs" -Filter "*.md" -Recurse -File | Where-Object { 
    $_.FullName -notlike "*\node_modules\*" 
}

Write-Host "📝 找到 $($mdFiles.Count) 个Markdown文件" -ForegroundColor Green

foreach ($file in $mdFiles) {
    $content = Get-Content -Path $file.FullName -Raw -Encoding UTF8
    $originalContent = $content
    $fileIssues = 0
    
    # 模式1: 移除连续重复的版本信息块
    # 匹配: **版本**: ... ---\n\n**版本**: ... ---
    $pattern1 = '(\*\*版本\*\*:.*?\n\*\*更新日期\*\*:.*?\n\*\*适用[^:]*\*\*:.*?\n+---\n+)\1+'
    if ($content -match $pattern1) {
        $before = $content
        $content = $content -replace $pattern1, '$1'
        if ($before -ne $content) {
            $fileIssues++
            Write-Host "  修复重复版本块: $($file.Name)" -ForegroundColor Yellow
        }
    }
    
    # 模式2: 移除多余的连续分隔符 (3个或更多)
    $pattern2 = '---\n+---\n+---'
    if ($content -match $pattern2) {
        $before = $content
        $content = $content -replace '---\n+---\n+---(\n+---)*', "---`n`n"
        if ($before -ne $content) {
            $fileIssues++
            Write-Host "  修复多余分隔符: $($file.Name)" -ForegroundColor Yellow
        }
    }
    
    # 模式3: 修复两个连续的 --- 为一个
    $pattern3 = '---\n+---\n+(?!#)'
    if ($content -match $pattern3) {
        $before = $content
        $content = $content -replace $pattern3, "---`n`n"
        if ($before -ne $content) {
            $fileIssues++
            Write-Host "  修复双分隔符: $($file.Name)" -ForegroundColor Yellow
        }
    }
    
    # 模式4: 特定模式 - 版本信息后紧跟两个 ---
    # **版本**: ...\n---\n\n---
    $pattern4 = '(\*\*版本\*\*:.*?\n\*\*更新日期\*\*:.*?\n\*\*适用[^:]*\*\*:.*?\n+)---\n+---\n+'
    if ($content -match $pattern4) {
        $before = $content
        $content = $content -replace $pattern4, "`$1---`n`n"
        if ($before -ne $content) {
            $fileIssues++
            Write-Host "  修复版本后分隔符: $($file.Name)" -ForegroundColor Yellow
        }
    }
    
    # 如果内容有变化,保存文件
    if ($originalContent -ne $content) {
        Set-Content -Path $file.FullName -Value $content -Encoding UTF8 -NoNewline
        $filesFixed++
        $totalIssues += $fileIssues
        Write-Host "✅ 已修复: $($file.Name) ($fileIssues 个问题)" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "=" * 60 -ForegroundColor Cyan
Write-Host "✨ 修复完成!" -ForegroundColor Green
Write-Host "📊 统计:" -ForegroundColor Cyan
Write-Host "  - 扫描文件: $($mdFiles.Count)" -ForegroundColor White
Write-Host "  - 修复文件: $filesFixed" -ForegroundColor Green
Write-Host "  - 修复问题: $totalIssues" -ForegroundColor Green
Write-Host "=" * 60 -ForegroundColor Cyan

