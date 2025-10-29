# 文档格式修复脚本
# 版本: 1.0
# 日期: 2025-10-29
# 用途: 批量修复docs目录下的格式问题

param(
    [string]$Path = "docs",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$stats = @{
    MetadataFixed = 0
    TitlesFixed = 0
    FilesProcessed = 0
    Errors = 0
}

Write-Host "🚀 开始文档格式修复..." -ForegroundColor Cyan
Write-Host "工作目录: $Path"
Write-Host "模式: $(if($DryRun){'试运行'}else{'实际修复'})`n"

#region 函数定义

function Fix-Metadata {
    param($FilePath)
    
    $content = Get-Content $FilePath -Raw -Encoding UTF8
    $modified = $false
    
    # 检测并替换"基准日期"格式
    if ($content -match '\*\*基准日期\*\*:') {
        if ($Verbose) { Write-Host "  [元数据] 发现旧格式: 基准日期" -ForegroundColor Yellow }
        
        # 提取Go版本
        $goVersion = "Go 1.25.3"
        if ($content -match '\*\*Go版本\*\*:\s*(.+?)(\r?\n)') {
            $goVersion = $matches[1].Trim()
        }
        
        # 构建新元数据
        $newMeta = @"
**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: $goVersion

---
"@
        
        # 替换（保留标题）
        $content = $content -replace '(?sm)(\r?\n)\*\*基准日期\*\*:.*?(\r?\n---)', "`n$newMeta"
        $modified = $true
    }
    
    # 统一日期格式: "2025年10月28日" -> "2025-10-29"
    if ($content -match '\*\*更新日期\*\*:\s*\d{4}年\d{1,2}月\d{1,2}日') {
        if ($Verbose) { Write-Host "  [元数据] 修复日期格式" -ForegroundColor Yellow }
        $content = $content -replace '\*\*更新日期\*\*:\s*\d{4}年\d{1,2}月\d{1,2}日', '**更新日期**: 2025-10-29'
        $modified = $true
    }
    
    if ($modified) {
        if (-not $DryRun) {
            Set-Content -Path $FilePath -Value $content -NoNewline -Encoding UTF8
        }
        $stats.MetadataFixed++
        return $true
    }
    
    return $false
}

function Fix-Title {
    param($FilePath)
    
    $content = Get-Content $FilePath -Raw -Encoding UTF8
    $modified = $false
    $fileName = Split-Path $FilePath -Leaf
    
    # 定义标题模式
    $titlePatterns = @{
        "00-概念定义体系.md" = @{
            Pattern = "^# .+ - 概念定义体系"
            Replace = "# 概念定义体系"
        }
        "00-知识图谱.md" = @{
            Pattern = "^# .+ - 知识图谱"
            Replace = "# 知识图谱"
        }
        "00-对比矩阵.md" = @{
            Pattern = "^# .+ - 对比矩阵"
            Replace = "# 对比矩阵"
        }
    }
    
    if ($titlePatterns.ContainsKey($fileName)) {
        $pattern = $titlePatterns[$fileName]
        
        if ($content -match $pattern.Pattern) {
            if ($Verbose) { Write-Host "  [标题] 修复: $fileName" -ForegroundColor Yellow }
            $content = $content -replace $pattern.Pattern, $pattern.Replace
            $modified = $true
        }
    }
    
    if ($modified) {
        if (-not $DryRun) {
            Set-Content -Path $FilePath -Value $content -NoNewline -Encoding UTF8
        }
        $stats.TitlesFixed++
        return $true
    }
    
    return $false
}

#endregion

#region 主执行逻辑

try {
    # 获取所有markdown文件
    $files = Get-ChildItem -Path $Path -Recurse -Filter "*.md" -File
    $totalFiles = $files.Count
    
    Write-Host "找到 $totalFiles 个Markdown文件`n" -ForegroundColor Green
    
    $progress = 0
    foreach ($file in $files) {
        $progress++
        $percentComplete = [math]::Round(($progress / $totalFiles) * 100)
        
        Write-Progress -Activity "处理文档" -Status "$progress/$totalFiles - $($file.Name)" -PercentComplete $percentComplete
        
        try {
            $stats.FilesProcessed++
            
            if ($Verbose) {
                Write-Host "处理 [$progress/$totalFiles]: $($file.FullName.Replace($PWD.Path, '.'))" -ForegroundColor Gray
            }
            
            # 执行修复
            $metaFixed = Fix-Metadata -FilePath $file.FullName
            $titleFixed = Fix-Title -FilePath $file.FullName
            
            if ($metaFixed -or $titleFixed) {
                Write-Host "  ✓ 修复: $($file.Name)" -ForegroundColor Green
            }
            
        } catch {
            $stats.Errors++
            Write-Host "  ✗ 错误: $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Progress -Activity "处理文档" -Completed
    
} catch {
    Write-Host "`n❌ 发生错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

#endregion

#region 结果报告

Write-Host "`n" + ("="*60) -ForegroundColor Cyan
Write-Host "📊 修复统计报告" -ForegroundColor Cyan
Write-Host ("="*60) -ForegroundColor Cyan

Write-Host "`n文件处理:"
Write-Host "  ✓ 已处理文件: $($stats.FilesProcessed)" -ForegroundColor Green
Write-Host "  ⚠ 错误数量: $($stats.Errors)" -ForegroundColor $(if($stats.Errors -gt 0){'Red'}else{'Gray'})

Write-Host "`n修复详情:"
Write-Host "  📝 元数据修复: $($stats.MetadataFixed) 个文件" -ForegroundColor Yellow
Write-Host "  🏷️  标题修复: $($stats.TitlesFixed) 个文件" -ForegroundColor Yellow
Write-Host "  📋 总修复: $($stats.MetadataFixed + $stats.TitlesFixed) 次" -ForegroundColor Yellow

if ($DryRun) {
    Write-Host "`n⚠️  这是试运行，未实际修改文件" -ForegroundColor Yellow
    Write-Host "执行实际修复请移除 -DryRun 参数" -ForegroundColor Yellow
} else {
    Write-Host "`n✅ 修复已完成！" -ForegroundColor Green
}

Write-Host "`n" + ("="*60) -ForegroundColor Cyan

#endregion

# 返回统计信息
return $stats

