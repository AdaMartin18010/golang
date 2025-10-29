# 文档格式修复脚本 v2.0
# 修复更多格式问题：重复目录、元数据变体、标题规范化

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    DuplicateTOC = 0
    MetadataFixed = 0
    TitlesFixed = 0
    FilesProcessed = 0
    Errors = 0
}

Write-Host "🚀 文档格式修复 v2.0..." -ForegroundColor Cyan
Write-Host "模式: $(if($DryRun){'试运行'}else{'实际修复'})`n"

function Remove-DuplicateTOC {
    param($Content)
    
    $modified = $false
    
    # 移除自动生成的TOC（保留手动的）
    if ($Content -match '(?sm)<!-- TOC START -->.*?<!-- TOC END -->\s*\r?\n') {
        if ($Verbose) { Write-Host "    [目录] 移除自动生成的TOC" -ForegroundColor Yellow }
        $Content = $Content -replace '(?sm)<!-- TOC START -->.*?<!-- TOC END -->\s*\r?\n', ''
        $modified = $true
        $stats.DuplicateTOC++
    }
    
    return @($Content, $modified)
}

function Fix-MetadataQuoteStyle {
    param($Content)
    
    $modified = $false
    
    # 转换引用格式的元数据为普通格式
    if ($Content -match '(?sm)^>\s*\*\*简介\*\*:.*?\r?\n>\s*\*\*版本\*\*:') {
        if ($Verbose) { Write-Host "    [元数据] 转换引用格式" -ForegroundColor Yellow }
        
        # 提取信息
        $version = "Go 1.25.3"
        if ($Content -match '>\s*\*\*版本\*\*:\s*(.+?)(?:\r?\n|$)') {
            $version = $matches[1].Trim()
        }
        
        # 替换为标准格式
        $newMeta = @"
**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: $version

---
"@
        
        # 移除旧的引用格式元数据
        $Content = $Content -replace '(?sm)^>\s*\*\*简介\*\*:.*?(?:\r?\n\r?\n|(?=##))', "$newMeta`n`n"
        $modified = $true
        $stats.MetadataFixed++
    }
    
    return @($Content, $modified)
}

function Fix-LongTitle {
    param($Content, $FileName)
    
    $modified = $false
    
    # 标题规范化规则
    $titleRules = @{
        "Go调度器与G-P-M模型" = "Go调度器"
        "Go并发编程进阶深度指南" = "Go并发编程进阶"
        "Go-1.25.3并发编程完整实战" = "并发编程完整实战"
    }
    
    foreach ($old in $titleRules.Keys) {
        if ($Content -match "^# $old\r?\n") {
            if ($Verbose) { Write-Host "    [标题] 简化: $old -> $($titleRules[$old])" -ForegroundColor Yellow }
            $Content = $Content -replace "^# $old\r?\n", "# $($titleRules[$old])`n"
            $modified = $true
            $stats.TitlesFixed++
            break
        }
    }
    
    return @($Content, $modified)
}

function Fix-TOCLink {
    param($Content)
    
    $modified = $false
    
    # 修复目录链接格式（移除多余的空格和特殊字符）
    # 例如: [1. 理论基础](#1-理论基础) 是正确的
    # 如果发现格式问题，修复它
    
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
        
        Write-Progress -Activity "处理文档 v2" -Status "$progress/$totalFiles - $($file.Name)" -PercentComplete $percentComplete
        
        try {
            $stats.FilesProcessed++
            $content = Get-Content $file.FullName -Raw -Encoding UTF8
            $originalContent = $content
            $hasChanges = $false
            
            if ($Verbose) {
                Write-Host "[$progress/$totalFiles] $($file.FullName.Replace($PWD.Path, '.'))" -ForegroundColor Gray
            }
            
            # 1. 移除重复TOC
            $result = Remove-DuplicateTOC -Content $content
            $content = $result[0]
            $hasChanges = $hasChanges -or $result[1]
            
            # 2. 修复引用格式的元数据
            $result = Fix-MetadataQuoteStyle -Content $content
            $content = $result[0]
            $hasChanges = $hasChanges -or $result[1]
            
            # 3. 简化过长标题
            $result = Fix-LongTitle -Content $content -FileName $file.Name
            $content = $result[0]
            $hasChanges = $hasChanges -or $result[1]
            
            # 保存修改
            if ($hasChanges -and -not $DryRun) {
                Set-Content -Path $file.FullName -Value $content -NoNewline -Encoding UTF8
                Write-Host "  ✓ 修复: $($file.Name)" -ForegroundColor Green
            }
            elseif ($hasChanges -and $DryRun) {
                Write-Host "  [DRY] 将修复: $($file.Name)" -ForegroundColor Yellow
            }
            
        } catch {
            $stats.Errors++
            Write-Host "  ✗ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Progress -Activity "处理文档 v2" -Completed
    
} catch {
    Write-Host "`n❌ 发生错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 结果报告
Write-Host "`n" + ("="*60) -ForegroundColor Cyan
Write-Host "📊 修复统计报告 v2.0" -ForegroundColor Cyan
Write-Host ("="*60) -ForegroundColor Cyan

Write-Host "`n文件处理:"
Write-Host "  ✓ 已处理文件: $($stats.FilesProcessed)" -ForegroundColor Green
Write-Host "  ⚠ 错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -gt 0){'Red'}else{'Gray'})

Write-Host "`n修复详情:"
Write-Host "  📑 移除重复TOC: $($stats.DuplicateTOC) 个文件" -ForegroundColor Yellow
Write-Host "  📝 元数据格式: $($stats.MetadataFixed) 个文件" -ForegroundColor Yellow
Write-Host "  🏷️  标题简化: $($stats.TitlesFixed) 个文件" -ForegroundColor Yellow
Write-Host "  📋 总修复: $($stats.DuplicateTOC + $stats.MetadataFixed + $stats.TitlesFixed) 次" -ForegroundColor Yellow

if ($DryRun) {
    Write-Host "`n⚠️  这是试运行，未实际修改文件" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复已完成！" -ForegroundColor Green
}

Write-Host "`n" + ("="*60) -ForegroundColor Cyan

return $stats

