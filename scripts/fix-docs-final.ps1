# 文档格式修复最终版
# 处理所有剩余的元数据变体

param(
    [string]$Path = "docs",
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"
$stats = @{
    MetadataFixed = 0
    FilesProcessed = 0
}

Write-Host "🎯 最终格式修复..." -ForegroundColor Cyan

function Normalize-AllMetadata {
    param($Content)
    
    $modified = $false
    
    # 策略：统一所有引用格式，保留简介/难度/标签在正文前
    if ($Content -match '(?sm)^>\s*\*\*') {
        $version = "Go 1.25.3"
        $date = "2025-10-29"
        $intro = ""
        $difficulty = ""
        $tags = ""
        
        # 提取所有字段
        if ($Content -match '>\s*\*\*(?:版本|Version)\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $version = $matches[1].Trim()
        }
        if ($Content -match '>\s*\*\*(?:更新日期|日期|Date)\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $date = $matches[1].Trim()
        }
        if ($Content -match '>\s*\*\*(?:适用于|Applies)\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $version = $matches[1].Trim()
        }
        if ($Content -match '>\s*\*\*简介\*\*:\s*(.+?)(?:\s*\r?\n)') {
            $intro = $matches[1].Trim()
        }
        if ($Content -match '>\s*\*\*难度\*\*:\s*(.+?)(?:\s*\r?\n)') {
            $difficulty = $matches[1].Trim()
        }
        if ($Content -match '>\s*\*\*标签\*\*:\s*(.+?)(?:\s*\r?\n)') {
            $tags = $matches[1].Trim()
        }
        
        # 构建新元数据
        $newMeta = @"
**版本**: v1.0  
**更新日期**: $date  
**适用于**: $version

---
"@
        
        # 如果有简介/难度/标签，放在前面
        $extras = ""
        if ($intro) { $extras += "> **简介**: $intro`n" }
        if ($difficulty) { $extras += "> **难度**: $difficulty`n" }
        if ($tags) { $extras += "> **标签**: $tags`n" }
        
        if ($extras) {
            $newMeta = "$extras`n$newMeta"
        }
        
        # 移除所有旧的引用格式元数据
        $Content = $Content -replace '(?sm)^>\s*\*\*[^:]+\*\*:.*?(?=\r?\n(?!>)|\r?\n\*\*|\r?\n##|\r?\n\r?\n---)', "$newMeta"
        
        $modified = $true
        $stats.MetadataFixed++
    }
    
    return @($Content, $modified)
}

try {
    $files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File | Where-Object { 
        $content = Get-Content $_.FullName -Raw
        $content -match '(?sm)^>\s*\*\*'
    }
    
    Write-Host "找到 $($files.Count) 个需要处理的文件`n" -ForegroundColor Yellow
    
    foreach ($file in $files) {
        $stats.FilesProcessed++
        $content = Get-Content $file.FullName -Raw -Encoding UTF8
        
        $result = Normalize-AllMetadata -Content $content
        $content = $result[0]
        $hasChanges = $result[1]
        
        if ($hasChanges -and -not $DryRun) {
            Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
            Write-Host "  ✓ $($file.Name)" -ForegroundColor Green
        }
    }
    
} catch {
    Write-Host "❌ 错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host "`n" + ("="*50) -ForegroundColor Cyan
Write-Host "📊 修复完成" -ForegroundColor Cyan
Write-Host ("="*50) -ForegroundColor Cyan
Write-Host "  处理文件: $($stats.FilesProcessed)" -ForegroundColor Green
Write-Host "  修复元数据: $($stats.MetadataFixed)" -ForegroundColor Green

return $stats

