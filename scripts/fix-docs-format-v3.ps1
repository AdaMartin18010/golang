# 文档格式修复脚本 v3.0
# 处理更多元数据变体和特殊情况

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    QuoteMetadata = 0
    SingleLineQuote = 0
    TitlesNormalized = 0
    FilesProcessed = 0
    Errors = 0
}

Write-Host "🚀 文档格式修复 v3.0 - 处理元数据变体..." -ForegroundColor Cyan
Write-Host "模式: $(if($DryRun){'试运行'}else{'实际修复'})`n"

function Fix-QuotedMetadata {
    param($Content, $FileName)
    
    $modified = $false
    
    # 匹配各种引用格式的元数据
    # 格式1: > **字段**: 值
    # 格式2: > **字段**: 值 (带换行)
    
    if ($Content -match '(?sm)^>\s*\*\*') {
        if ($Verbose) { Write-Host "    [元数据] 发现引用格式，转换中..." -ForegroundColor Yellow }
        
        # 提取各个字段
        $version = "Go 1.25.3"
        $date = "2025-10-29"
        $difficulty = ""
        $tags = ""
        $intro = ""
        
        # 版本
        if ($Content -match '>\s*\*\*版本\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $version = $matches[1].Trim()
        }
        
        # 更新日期（可能没有）
        if ($Content -match '>\s*\*\*更新日期\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $date = $matches[1].Trim()
        }
        elseif ($Content -match '>\s*\*\*日期\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $date = $matches[1].Trim()
        }
        
        # 简介（可能有）
        if ($Content -match '>\s*\*\*简介\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $intro = $matches[1].Trim()
        }
        
        # 难度（可能有）
        if ($Content -match '>\s*\*\*难度\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $difficulty = $matches[1].Trim()
        }
        
        # 标签（可能有）
        if ($Content -match '>\s*\*\*标签\*\*:\s*(.+?)(?:\s*\r?\n|$)') {
            $tags = $matches[1].Trim()
        }
        
        # 构建新的元数据
        $newMeta = @"
**版本**: v1.0  
**更新日期**: $date  
**适用于**: $version

---
"@
        
        # 如果有简介、难度、标签，添加到新格式中
        if ($intro -or $difficulty -or $tags) {
            $extras = ""
            if ($intro) { $extras += "> **简介**: $intro`n" }
            if ($difficulty) { $extras += "> **难度**: $difficulty`n" }
            if ($tags) { $extras += "> **标签**: $tags`n" }
            $newMeta = "$extras`n$newMeta"
        }
        
        # 移除整个引用块（从第一个 > ** 到下一个非引用行）
        $Content = $Content -replace '(?sm)^>\s*\*\*(?:简介|版本|更新日期|日期|难度|标签|适用于)\*\*:.*?(?:\r?\n(?!>)|\r?\n\r?\n)', "$newMeta`n`n"
        
        $modified = $true
        $stats.QuoteMetadata++
    }
    
    return @($Content, $modified)
}

function Normalize-FileTitle {
    param($Content, $FileName)
    
    $modified = $false
    
    # 从文件名提取预期标题
    # 例如: "04-Go调度器.md" -> "Go调度器"
    if ($FileName -match '^\d+-(.+?)\.md$') {
        $expectedTitle = $matches[1]
        
        # 检查实际标题是否过长或不一致
        if ($Content -match '^# (.+?)\r?\n') {
            $actualTitle = $matches[1]
            
            # 如果标题包含额外信息，简化它
            $simplifyRules = @{
                "(.+?)与.+-P-M模型" = '$1'
                "(.+?)进阶深度指南" = '$1进阶'
                "(.+?)深度实战指南" = '$1'
                "(.+?)完整实战指南" = '$1'
                "Go-1\.25\.3(.+?)完整实战" = '$1'
                "(.+?)-完整实现指南" = '$1'
            }
            
            $newTitle = $actualTitle
            foreach ($pattern in $simplifyRules.Keys) {
                if ($actualTitle -match $pattern) {
                    $newTitle = $actualTitle -replace $pattern, $simplifyRules[$pattern]
                    break
                }
            }
            
            if ($newTitle -ne $actualTitle) {
                if ($Verbose) { Write-Host "    [标题] 规范化: $actualTitle -> $newTitle" -ForegroundColor Yellow }
                $Content = $Content -replace "^# $([regex]::Escape($actualTitle))\r?\n", "# $newTitle`n"
                $modified = $true
                $stats.TitlesNormalized++
            }
        }
    }
    
    return @($Content, $modified)
}

function Fix-SingleLineQuoteMeta {
    param($Content)
    
    $modified = $false
    
    # 处理单行引用格式（在文档末尾或中间）
    # 例如: > 版本: v1.0 | 更新: 2025-10-29
    
    return @($Content, $modified)
}

try {
    $files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File
    $totalFiles = $files.Count
    
    Write-Host "找到 $totalFiles 个Markdown文件`n" -ForegroundColor Green
    
    $progress = 0
    foreach ($file in $files) {
        $progress++
        $percentComplete = [math]::Round(($progress / $totalFiles) * 100)
        
        Write-Progress -Activity "处理文档 v3" -Status "$progress/$totalFiles - $($file.Name)" -PercentComplete $percentComplete
        
        try {
            $stats.FilesProcessed++
            $content = Get-Content $file.FullName -Raw -Encoding UTF8
            $originalContent = $content
            $hasChanges = $false
            
            if ($Verbose -and ($content -match '(?sm)^>\s*\*\*')) {
                Write-Host "[$progress/$totalFiles] $($file.FullName.Replace($PWD.Path, '.'))" -ForegroundColor Gray
            }
            
            # 1. 修复引用格式的元数据
            $result = Fix-QuotedMetadata -Content $content -FileName $file.Name
            $content = $result[0]
            $hasChanges = $hasChanges -or $result[1]
            
            # 2. 规范化标题
            $result = Normalize-FileTitle -Content $content -FileName $file.Name
            $content = $result[0]
            $hasChanges = $hasChanges -or $result[1]
            
            # 保存修改
            if ($hasChanges -and -not $DryRun) {
                Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
                if (-not $Verbose) {
                    Write-Host "  ✓ 修复: $($file.Name)" -ForegroundColor Green
                }
            }
            elseif ($hasChanges -and $DryRun) {
                Write-Host "  [DRY] 将修复: $($file.Name)" -ForegroundColor Yellow
            }
            
        } catch {
            $stats.Errors++
            Write-Host "  ✗ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Progress -Activity "处理文档 v3" -Completed
    
} catch {
    Write-Host "`n❌ 发生错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 结果报告
Write-Host "`n" + ("="*60) -ForegroundColor Cyan
Write-Host "📊 修复统计报告 v3.0" -ForegroundColor Cyan
Write-Host ("="*60) -ForegroundColor Cyan

Write-Host "`n文件处理:"
Write-Host "  ✓ 已处理文件: $($stats.FilesProcessed)" -ForegroundColor Green
Write-Host "  ⚠ 错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -gt 0){'Red'}else{'Gray'})

Write-Host "`n修复详情:"
Write-Host "  📝 引用格式元数据: $($stats.QuoteMetadata) 个文件" -ForegroundColor Yellow
Write-Host "  🏷️  标题规范化: $($stats.TitlesNormalized) 个文件" -ForegroundColor Yellow
Write-Host "  📋 总修复: $($stats.QuoteMetadata + $stats.TitlesNormalized) 次" -ForegroundColor Yellow

if ($DryRun) {
    Write-Host "`n⚠️  这是试运行，未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复已完成！" -ForegroundColor Green
}

Write-Host "`n" + ("="*60) -ForegroundColor Cyan

return $stats

